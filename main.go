package main

import (
	"log"

	"diamond/handlers"
	"diamond/sshd"
)

func main() {
	// 启动sshd服务器
	go sshd.Start()
	// 启动web服务器
	app := handlers.App()
	log.Fatal(app.Listen(":3000"))
}
