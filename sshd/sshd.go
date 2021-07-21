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
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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
	// 连接数据库获取信息
	dsn := "root:#yz2NRz30d_>m^:n90^V@tcp(192.168.241.130:3306)/Bumblebee-new"
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Printf("connect server failed, err:%v\n", err)
		s.Exit(1)
	}
	servers := []Server{}
	db.Select(&servers, "SELECT id, ip, user, port, password FROM cmdb_server")
	data := [][]string{}
	for _, server := range servers {
		data = append(data, []string{strconv.Itoa(server.ID), server.IP, server.USER, strconv.Itoa(server.PORT), server.PASSWORD})
	}
	// 展示table信息
	table := tablewriter.NewWriter(s)
	table.SetHeader([]string{"id", "ip", "user", "port", "password"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
	// 开启ternimal跟用户交互
	promt := fmt.Sprintf("[%s]=> ", s.User())
	for {
		terminal := term.NewTerminal(s, "")
		terminal.SetPrompt(string(terminal.Escape.Red) + promt + string(terminal.Escape.Reset))
		io.WriteString(s, "please choose a server: ")
		line, err := terminal.ReadLine()
		if err == io.EOF {
			log.Println(err)
			io.WriteString(s, "\n")
			continue
		}
		if err != nil {
			log.Println(err)
			s.Exit(1)
			return
		}
		if line == "" {
			log.Println("empty")
			continue
		}
		if len(line) > 0 {
			io.WriteString(s, fmt.Sprintf("\nyou choose: %s \n", line))
			break
		}
	}

	// 连接远程服务器
	_, winCh, isPty := s.Pty()
	if isPty {
		client, err := utils.GetSSHClient("192.168.241.130", 22, "root", "12345678")
		if err != nil {
			log.Println(err)
			s.Exit(1)
		}
		defer client.Close()
		//创建ssh session
		session, err := client.NewSession()
		if err != nil {
			log.Printf("create session failed: %v", err)
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
			s.Exit(1)
		}
		// window change
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
