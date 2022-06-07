package sshd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"diamond/config"
	"diamond/models"
	"diamond/policy"
	"diamond/utils"

	"github.com/gliderlabs/ssh"
	"github.com/olekukonko/tablewriter"
	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

var (
	DeadlineTimeout = 1 * time.Hour
	IdleTimeout     = 30 * time.Minute
)

func passwordHandler(ctx ssh.Context, password string) bool {
	user := &models.User{}
	if result := models.DB.Where("username = ?", ctx.User()).First(user); result.Error != nil {
		return false
	}
	if !utils.CheckPassword(user.Password, password) {
		return false
	}
	if !user.IsActive {
		return false
	}
	return true
}

func sshHandler(s ssh.Session) {
	// 登录成功后提示信息
	io.WriteString(s, "\n***** Welcome to diamond *****\n")
	io.WriteString(s, fmt.Sprintf("\n***** Connections will only last %s *****\n", DeadlineTimeout))
	io.WriteString(s, fmt.Sprintf("\n***** Timeout after %s of no activity *****\n\n", IdleTimeout))
	// 输出服务器分组信息
	user := &models.User{}
	models.DB.Where("username = ?", s.User()).First(user)
	groupTable := tablewriter.NewWriter(s)
	groupTable.SetHeader([]string{"id", "name"})
	groupTable.Append([]string{"0", "ALL(所有服务器)"})
	groupMap := map[uint]string{}
	var gids []string // 普通用户服务器组
	if user.IsSuperuser {
		groups := &models.Groups{}
		models.DB.Find(groups)
		for _, v := range *groups {
			groupMap[v.ID] = v.Name
		}
	} else {
		sub := fmt.Sprintf("user::%d", user.ID)
		perms, _ := policy.Enforcer.GetNamedImplicitPermissionsForUser("p", sub)
		for _, perm := range perms {
			if utils.FindValInSlice(perm, "server") {
				gids = append(gids, strings.ReplaceAll(perm[1], "group::", ""))
			}
		}
		groups := &models.Groups{}
		models.DB.Where("id IN ?", gids).Find(&groups)
		for _, g := range *groups {
			if g.IsActive {
				groupMap[g.ID] = g.Name
			}
		}
	}
	for k, v := range groupMap {
		groupTable.Append([]string{strconv.Itoa(int(k)), v})
	}
	groupTable.Render()
	// 开启ternimal跟用户交互
	promt := fmt.Sprintf("[%s]=> ", s.User())
	terminal := term.NewTerminal(s, "")
	terminal.SetPrompt(string(terminal.Escape.Red) + promt + string(terminal.Escape.Reset))
	var groupInput string
	for {
		io.WriteString(s, "请选择分组ID: ")
		line, err := terminal.ReadLine()
		if err == io.EOF {
			log.Println(err)
			s.Exit(1)
			return
		}
		if err != nil {
			log.Println(err)
			s.Exit(1)
			return
		}
		if line == "" {
			log.Println("empty")
			groupTable.Render()
			continue
		}
		if len(line) > 0 {
			groupInput = line
			break
		}
	}
	groupID, err := strconv.Atoi(groupInput)
	if err != nil {
		io.WriteString(s, err.Error())
		s.Exit(1)
	}
	// 服务器查找
	serverMap := map[uint]models.Server{}
	rules := policy.Enforcer.GetNamedGroupingPolicy("g2")
	if groupID == 0 {
		log.Println("展示所有服务器")
		servers := &models.Servers{}
		if user.IsSuperuser {
			models.DB.Find(servers)
			for _, s := range *servers {
				serverMap[s.ID] = s
			}
		} else {
			sids := []string{}
			for _, rule := range rules {
				sids = append(sids, strings.ReplaceAll(rule[0], "server::", ""))
			}
			models.DB.Where("id IN ?", sids).Find(&servers)
			for _, server := range *servers {
				if server.IsActive {
					serverMap[server.ID] = server
				}
			}
		}
	} else {
		if _, ok := groupMap[uint(groupID)]; ok {
			log.Println("展示选择的组里的服务器")
			sids := []string{}
			for _, rule := range rules {
				x := fmt.Sprintf("group::%s", groupInput)
				if rule[1] == x {
					sids = append(sids, strings.ReplaceAll(rule[0], "server::", ""))
				}
			}
			servers := &models.Servers{}
			models.DB.Where("id IN ?", sids).Find(&servers)
			for _, s := range *servers {
				if user.IsSuperuser {
					serverMap[s.ID] = s
				} else {
					if s.IsActive {
						serverMap[s.ID] = s
					}
				}
			}
		} else {
			io.WriteString(s, "分组不存在！\n")
			s.Exit(1)
		}
	}
	serverTable := tablewriter.NewWriter(s)
	serverTable.SetHeader([]string{"id", "ip", "hostname", "remark"})
	for k, v := range serverMap {
		serverTable.Append([]string{strconv.Itoa(int(k)), v.IP, v.Hostname.String, v.Remark.String})
	}
	serverTable.Render()
	var serverInput string
	for {
		io.WriteString(s, "请选择服务器ID: ")
		line, err := terminal.ReadLine()
		if err == io.EOF {
			log.Println(err)
			s.Exit(1)
			return
		}
		if err != nil {
			log.Println(err)
			s.Exit(1)
			return
		}
		if line == "" {
			log.Println("empty")
			serverTable.Render()
			continue
		}
		if len(line) > 0 {
			serverInput = line
			break
		}
	}
	serverID, err := strconv.Atoi(serverInput)
	if err != nil {
		io.WriteString(s, err.Error())
		s.Exit(1)
	}
	if server, ok := serverMap[uint(serverID)]; !ok {
		io.WriteString(s, "服务器不存在！\n")
		s.Exit(1)
	} else {
		log.Printf("连接服务器%d...", server.ID)
		key := &models.Key{}
		models.DB.First(key, server.KeyID.UInt32)
		// 连接远程服务器
		_, winCh, isPty := s.Pty()
		if isPty {
			client, err := utils.GetSSHClient(server.IP, int(server.Port), server.User, int(server.AuthType), server.Password.String, key.Content)
			if err != nil {
				log.Println(err)
				io.WriteString(s, fmt.Sprintf("连接服务器失败: %v", err))
				s.Exit(1)
			}
			defer client.Close()
			//创建ssh session
			session, err := client.NewSession()
			if err != nil {
				log.Printf("create session failed: %v", err)
				io.WriteString(s, fmt.Sprintf("创建session失败: %v", err))
				s.Exit(1)
			}
			defer session.Close()
			// Set up terminal modes
			modes := gossh.TerminalModes{
				gossh.ECHO:          1,     // disable echoing
				gossh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
				gossh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
			}
			// Request pseudo terminal
			if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
				log.Printf("request for pseudo terminal failed: %v", err)
				io.WriteString(s, fmt.Sprintf("request for pseudo terminal failed: %v", err))
				s.Exit(1)
			}
			session.Stdin = s
			// session.Stdout = s
			session.Stderr = s
			// 开始记录终端录屏
			now := time.Now()
			time_str := now.Format("2006-01-02-15-04-05")
			record_dir := "./records/"
			if !utils.PathExists(record_dir) {
				os.Mkdir(record_dir, 0755)
			}
			record_file := fmt.Sprintf("%s%s-%d-%s.cast", record_dir, s.User(), server.ID, time_str)
			record := models.Record{User: s.User(), ServerID: server.ID, File: record_file}
			models.DB.Create(&record)
			// 捕捉输出，输出有回显，可用来记录终端输入输出
			out, _ := session.StdoutPipe()
			go func() {
				var f *os.File
				f, _ = os.Create(record_file)
				defer f.Close()
				// 记录文件头
				timestamp, _ := strconv.ParseFloat(fmt.Sprintf("%.9f", float64(now.UnixNano())/float64(1e9)), 64)
				env := map[string]string{
					"SHELL": "/bin/bash",
					"TERM":  "xterm",
				}
				header := map[string]interface{}{
					"version":   2,
					"width":     80,
					"height":    40,
					"timestamp": timestamp,
					"env":       env,
				}
				headerbyte, _ := json.Marshal(header)
				f.WriteString(string(headerbyte) + "\n")
				// 记录终端内容
				for {
					var buffer [1024]byte
					n, err := out.Read(buffer[:])
					if err != nil {
						fmt.Println("out error:", err)
						break
					}
					nt, _ := strconv.ParseFloat(fmt.Sprintf("%.9f", float64(time.Now().UnixNano())/float64(1e9)), 64)
					iodata := []string{fmt.Sprintf("%.9f", nt-timestamp), "o", string(buffer[:n])}
					iodatabyte, _ := json.Marshal(iodata)
					f.WriteString(string(iodatabyte) + "\n")
					// fmt.Println(string(buffer[:n]))
					s.Write(buffer[:n])
				}
			}()
			// Start remote shell
			if err := session.Shell(); err != nil {
				log.Printf("failed to start shell: %v", err)
				io.WriteString(s, fmt.Sprintf("failed to start shell: %v", err))
				s.Exit(1)
			}
			// 监听window change
			go func() {
				for win := range winCh {
					log.Printf("change window: %d %d", win.Height, win.Width)
					session.WindowChange(win.Height, win.Width)
				}
			}()
			// 监听超时, 关闭session
			go func() {
				for {
					select {
					case <-time.After(time.Second):
						continue
					case <-s.Context().Done():
						log.Println("connection closed")
						session.Close()
						s.Exit(1)
						return
					}
				}
			}()
			session.Wait()
		} else {
			io.WriteString(s, "No PTY requested.\n")
			s.Exit(1)
		}
	}
}

func Start() {
	addr := config.Config.Get("sshd.addr").(string)
	keyPath := config.Config.Get("sshd.keyPath").(string)
	server := &ssh.Server{
		Addr:            addr,
		MaxTimeout:      DeadlineTimeout,
		IdleTimeout:     IdleTimeout,
		PasswordHandler: passwordHandler,
		Handler:         sshHandler,
	}
	dat, _ := ioutil.ReadFile(keyPath)
	key, _ := gossh.ParsePrivateKey(dat)
	server.AddHostKey(key)
	log.Printf("starting sshd server on %s", addr)
	log.Fatal(server.ListenAndServe())
}
