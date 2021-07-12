package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

var (
	DeadlineTimeout = 6 * time.Hour
	IdleTimeout     = 30 * time.Minute
)

func passwordHandler(ctx ssh.Context, password string) bool {
	log.Println(ctx.SessionID())
	return password == "secret"
}

func sshHandler(s ssh.Session) {
	// io.WriteString(s, fmt.Sprintf("Hello %s\n", s.User()))
	_, winCh, isPty := s.Pty()
	if isPty {
		sshPort := 22
		sshHost := "13.212.227.122"
		sshUser := "root"
		sshPassword := "LpNkH&#Y33sbKfDlMCFR"

		config := &gossh.ClientConfig{
			Timeout:         time.Second,
			User:            sshUser,
			HostKeyCallback: gossh.InsecureIgnoreHostKey(),
			Auth:            []gossh.AuthMethod{gossh.Password(sshPassword)},
		}
		addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
		sshClient, err := gossh.Dial("tcp", addr, config)
		if err != nil {
			log.Fatalf("创建ssh client 失败: %v", err.Error())
		}
		defer sshClient.Close()
		//创建ssh-session
		session, err := sshClient.NewSession()
		if err != nil {
			log.Fatalf("创建ssh session 失败: %v", err.Error())
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
			log.Fatal("request for pseudo terminal failed: ", err)
		}
		// session.Stdin = s
		// session.Stdout = s
		session.Stderr = s
		// 捕捉输入，可用来过滤用户输入
		in, _ := session.StdinPipe()
		go func() {
			for {
				var buffer [1024]byte
				n, err := s.Read(buffer[:])
				if err != nil {
					fmt.Println("in error:", err)
					break
				}
				fmt.Println(string(buffer[:n]))
				in.Write(buffer[:n])
			}
		}()
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
			log.Fatal("failed to start shell: ", err)
		}
		// window change
		go func() {
			for win := range winCh {
				log.Printf("change window: %d %d", win.Height, win.Width)
				session.WindowChange(win.Height, win.Width)
			}
		}()
		session.Wait()
	} else {
		io.WriteString(s, "No PTY requested.\n")
		s.Exit(1)
	}
}

func main() {
	log.Println("starting ssh server on port 2222...")
	log.Printf("connections will only last %s\n", DeadlineTimeout)
	log.Printf("and timeout after %s of no activity\n", IdleTimeout)

	server := &ssh.Server{
		Addr:            ":2222",
		MaxTimeout:      DeadlineTimeout,
		IdleTimeout:     IdleTimeout,
		PasswordHandler: passwordHandler,
		Handler:         sshHandler,
	}
	dat, _ := ioutil.ReadFile("C:/Users/90hua/.ssh/id_rsa")
	key, _ := gossh.ParsePrivateKey(dat)
	server.AddHostKey(key)
	log.Fatal(server.ListenAndServe())
}
