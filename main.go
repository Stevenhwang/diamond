package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	"diamond/utils.go"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

var (
	DeadlineTimeout = 6 * time.Hour
	IdleTimeout     = 10 * time.Second
)

func passwordHandler(ctx ssh.Context, password string) bool {
	log.Println(ctx.SessionID())
	return password == "secret"
}

func sshHandler(s ssh.Session) {
	io.WriteString(s, fmt.Sprintf("connections will only last %s\n", DeadlineTimeout))
	io.WriteString(s, fmt.Sprintf("and timeout after %s of no activity\n", IdleTimeout))

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
		session.Wait()
	} else {
		io.WriteString(s, "No PTY requested.\n")
		s.Exit(1)
	}
}

func main() {
	log.Println("starting ssh server on port 2222...")

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
