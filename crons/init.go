package crons

import (
	"log"

	"github.com/robfig/cron/v3"
)

// 初始化定时任务
func init() {
	c := cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger)))
	log.Println("开启定时任务")
	c.AddFunc("@daily", CleanLogTask)
	c.Start()
}
