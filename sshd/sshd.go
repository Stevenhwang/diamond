package sshd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"diamond/config"
	"diamond/models"
	"diamond/utils.go"

	"github.com/gliderlabs/ssh"
	"github.com/olekukonko/tablewriter"
	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

var (
	DeadlineTimeout = 3 * time.Hour
	IdleTimeout     = 30 * time.Minute
)

type Server struct {
	ID       int    `db:"id"`
	IP       string `db:"ip"`
	USER     string `db:"user"`
	PORT     int    `db:"port"`
	PASSWORD string `db:"password"`
}

func passwordHandler(ctx ssh.Context, password string) bool {
	// return password == "secret"
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
	models.DB.Preload("Roles.Groups").Where("username = ?", s.User()).First(user)
	groupTable := tablewriter.NewWriter(s)
	groupTable.SetHeader([]string{"id", "name"})
	groupMap := map[uint]string{}
	if user.IsSuperuser {
		groups := &models.Groups{}
		models.DB.Find(groups)
		for _, v := range *groups {
			groupMap[v.ID] = v.Name
		}
	} else {
		for _, role := range user.Roles {
			if role.IsActive {
				for _, group := range role.Groups {
					if group.IsActive {
						groupMap[group.ID] = group.Name
					}
				}
			}
		}
	}
	gIDList := []int{}
	for k, v := range groupMap {
		gIDList = append(gIDList, int(k))
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
	serverMap := map[uint]models.Server{}
	if groupID == 0 {
		log.Println("展示所有服务器")
		if user.IsSuperuser {
			servers := &models.Servers{}
			models.DB.Find(servers)
			for _, s := range *servers {
				serverMap[s.ID] = s
			}
		} else {
			groups := &models.Groups{}
			models.DB.Preload("Servers").Find(groups, gIDList)
			for _, g := range *groups {
				for _, server := range g.Servers {
					if server.IsActive {
						serverMap[server.ID] = server
					}
				}
			}
		}
	} else {
		if _, ok := groupMap[uint(groupID)]; ok {
			log.Println("展示选择的组里的服务器")
			group := &models.Group{}
			models.DB.Preload("Servers").Find(group, groupID)
			if user.IsSuperuser {
				for _, s := range group.Servers {
					serverMap[s.ID] = s
				}
			} else {
				for _, s := range group.Servers {
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
	serverTable.SetHeader([]string{"id", "ip", "remark", "user", "port"})
	for k, v := range serverMap {
		serverTable.Append([]string{strconv.Itoa(int(k)), v.IP, v.Remark.String, v.User, strconv.Itoa(v.Port)})
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
		// 连接远程服务器
		_, winCh, isPty := s.Pty()
		if isPty {
			client, err := utils.GetSSHClient(server.IP, server.Port, server.User, server.Password.String, server.Key.String)
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
			// 捕捉输出，输出有回显，可用来记录终端输入输出
			out, _ := session.StdoutPipe()
			go func() {
				for {
					var buffer [1024]byte
					n, err := out.Read(buffer[:])
					if err != nil {
						fmt.Println("out error:", err)
						break
					}
					fmt.Println(string(buffer[:n]))
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
	log.Println("starting sshd server ...")
	log.Fatal(server.ListenAndServe())
}
