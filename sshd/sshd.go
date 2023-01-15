package sshd

import (
	"diamond/cache"
	"diamond/misc"
	"diamond/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	myssh "github.com/Stevenhwang/gommon/ssh"
	"github.com/Stevenhwang/gommon/times"
	"github.com/Stevenhwang/gommon/tools"
	"github.com/gliderlabs/ssh"
	"github.com/olekukonko/tablewriter"
	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/term"
	"gorm.io/gorm"
)

var (
	DeadlineTimeout = 30 * time.Minute
	IdleTimeout     = 5 * time.Minute
)

func passwordHandler(ctx ssh.Context, password string) bool {
	// 先检查IP是否在黑名单，再继续
	ip := ctx.RemoteAddr().(*net.TCPAddr).IP.String()
	b, err := cache.GetBan(ip)
	if err != nil {
		return false
	}
	if b {
		return false
	}
	user := models.User{}
	if res := models.DB.Where("username = ?", ctx.User()).First(&user); res.Error != nil {
		cache.Ban(ip) // 试错也加入黑名单
		return false
	}
	if !tools.CheckPassword(user.Password, password) {
		cache.Ban(ip)
		return false
	}
	if !user.IsActive {
		return false
	}
	return true
}

func publickeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	// 先检查IP是否在黑名单，再继续
	ip := ctx.RemoteAddr().(*net.TCPAddr).IP.String()
	b, err := cache.GetBan(ip)
	if err != nil {
		return false
	}
	if b {
		return false
	}
	user := models.User{}
	if res := models.DB.Where("username = ?", ctx.User()).First(&user); res.Error != nil {
		cache.Ban(ip) // 试错也加入黑名单
		return false
	}
	data := user.Publickey
	if len(data) == 0 { // 如果用户没有绑定key
		return false
	}
	allowed, _, _, _, _ := ssh.ParseAuthorizedKey([]byte(data))
	if !ssh.KeysEqual(key, allowed) { // key不匹配
		cache.Ban(ip)
		return false
	}
	return true
}

func sshHandler(s ssh.Session) {
	// 登录成功后提示信息
	io.WriteString(s, "================================\r\n")
	io.WriteString(s, fmt.Sprintf("连接最长%s, 空闲超时%s\r\n", DeadlineTimeout, IdleTimeout))
	io.WriteString(s, "================================\r\n")
	// 获取当前用户的服务器信息
	user := models.User{}
	servers := models.Servers{}
	if s.User() == "admin" {
		if res := models.DB.Order("created_at desc").Find(&servers); res.Error != nil {
			io.WriteString(s, res.Error.Error())
			s.Exit(1)
			return
		}
	} else {
		if res := models.DB.Where("username = ?", s.User()).Preload("Servers", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc")
		}).First(&user); res.Error != nil {
			io.WriteString(s, res.Error.Error())
			s.Exit(1)
			return
		}
		servers = user.Servers
	}
	// 服务器查找
	serverMap := map[int]models.Server{}
	for i, s := range servers {
		serverMap[i+1] = s
	}
	serverData := [][]string{}
	for k, v := range serverMap {
		serverData = append(serverData, []string{strconv.Itoa(k), v.Name, v.IP, v.Remark})
	}
	serverTable := tablewriter.NewWriter(s)
	serverTable.SetHeader([]string{"id", "name", "ip", "remark"})
	serverTable.AppendBulk(serverData)
	serverTable.Render()
	// 开启ternimal跟用户交互
	promt := fmt.Sprintf("[%s]=> ", s.User())
	terminal := term.NewTerminal(s, "")
	terminal.SetPrompt(string(terminal.Escape.Red) + promt + string(terminal.Escape.Reset))

	var serverInput string
	for {
		io.WriteString(s, "请选择服务器ID: ")
		line, err := terminal.ReadLine()
		if err == io.EOF {
			misc.Logger.Error().Err(err).Str("from", "sshd").Msg("")
			s.Exit(1)
			return
		}
		if err != nil {
			misc.Logger.Error().Err(err).Str("from", "sshd").Msg("")
			s.Exit(1)
			return
		}
		if line == "" {
			misc.Logger.Error().Err(err).Str("from", "sshd").Msg("empty")
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
		return
	}
	if server, ok := serverMap[serverID]; !ok {
		io.WriteString(s, "服务器不存在！\n")
		s.Exit(1)
		return
	} else {
		misc.Logger.Info().Str("from", "sshd").Msg("连接服务器")
		// 开始连接流程
		credential := models.Credential{}
		if res := models.DB.First(&credential, server.CredentialID); res.Error != nil {
			io.WriteString(s, "认证信息不存在")
			s.Exit(1)
			return
		}
		var client *gossh.Client
		spass := tools.AesDecrypt(credential.AuthContent, "0123456789012345")
		if credential.AuthType == 1 {
			cli, err := myssh.GetSSHClientByPassword(server.IP, spass, myssh.SSHOptions{Port: server.Port, User: credential.AuthUser})
			if err != nil {
				io.WriteString(s, err.Error())
				s.Exit(1)
				return
			}
			client = cli
		}
		if credential.AuthType == 2 {
			cli, err := myssh.GetSSHClientByKey(server.IP, []byte(spass), myssh.SSHOptions{Port: server.Port, User: credential.AuthUser})
			if err != nil {
				io.WriteString(s, err.Error())
				s.Exit(1)
				return
			}
			client = cli
		}
		defer client.Close()
		// 连接远程服务器
		_, winCh, isPty := s.Pty()
		if isPty {
			//创建ssh session
			session, err := client.NewSession()
			if err != nil {
				io.WriteString(s, fmt.Sprintf("创建session失败: %v", err))
				s.Exit(1)
				return
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
				io.WriteString(s, fmt.Sprintf("request for pseudo terminal failed: %v", err))
				s.Exit(1)
			}
			session.Stdin = s
			// session.Stdout = s
			session.Stderr = s
			// 开始记录终端录屏
			now := time.Now()
			// time_str := now.Format("2006-01-02-15-04-05")
			time_str, _ := times.GetTimeString(0, times.TimeOptions{Layout: "20060102150405"})
			record_dir := "./records/"
			record_file := fmt.Sprintf("%s%s-%s-%s-sshd.cast", record_dir, s.User(), server.IP, time_str)
			fromIP := s.RemoteAddr().(*net.TCPAddr).IP.String()
			record := models.Record{User: s.User(), IP: server.IP, FromIP: fromIP, File: record_file}
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
					"height":    24,
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
						misc.Logger.Error().Err(err).Str("from", "sshd").Msg("")
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
				io.WriteString(s, fmt.Sprintf("failed to start shell: %v", err))
				s.Exit(1)
				return
			}
			// 监听window change
			go func() {
				for win := range winCh {
					misc.Logger.Info().Str("from", "sshd").Msg("window change")
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
						misc.Logger.Info().Str("from", "sshd").Msg("connection closed")
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

func Start(addr string, keyPath string) {
	forwardHandler := &ssh.ForwardedTCPHandler{}
	server := &ssh.Server{
		Addr:             addr,
		MaxTimeout:       DeadlineTimeout,
		IdleTimeout:      IdleTimeout,
		PasswordHandler:  passwordHandler,
		PublicKeyHandler: publickeyHandler,
		Handler:          sshHandler,
		LocalPortForwardingCallback: ssh.LocalPortForwardingCallback(func(ctx ssh.Context, dhost string, dport uint32) bool {
			misc.Logger.Info().Str("from", "sshd").Str("user", ctx.User()).Msg(fmt.Sprintf("Accepted forward host:%s port:%d", dhost, dport))
			return true
		}),
		ReversePortForwardingCallback: ssh.ReversePortForwardingCallback(func(ctx ssh.Context, host string, port uint32) bool {
			misc.Logger.Info().Str("from", "sshd").Str("user", ctx.User()).Msg(fmt.Sprintf("Attempt to bind host:%s port:%d", host, port))
			return true
		}),
		RequestHandlers: map[string]ssh.RequestHandler{
			"tcpip-forward":        forwardHandler.HandleSSHRequest,
			"cancel-tcpip-forward": forwardHandler.HandleSSHRequest,
		},
		ChannelHandlers: map[string]ssh.ChannelHandler{
			"session":      ssh.DefaultSessionHandler,
			"direct-tcpip": ssh.DirectTCPIPHandler,
		},
	}
	dat, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatal(err)
	}
	key, err := gossh.ParsePrivateKey(dat)
	if err != nil {
		log.Fatal(err)
	}
	server.AddHostKey(key)
	misc.Logger.Info().Str("from", "sshd").Msg(fmt.Sprintf("starting sshd server on %s", addr))
	log.Fatal(server.ListenAndServe())
}
