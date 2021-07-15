package main

import (
	"log"

	"diamond/config"
	"diamond/handlers"
	"diamond/sshd"
)

func main() {
	// 启动sshd服务器
	go sshd.Start()
	// 启动web服务器
	app := handlers.App()
	addr := config.Config.Get("web.addr").(string)
	log.Fatal(app.Listen(addr))
}
