package crons

import (
	"log"

	"github.com/robfig/cron/v3"
)

var C *cron.Cron

// 初始化定时任务
func init() {
	C = cron.New(cron.WithSeconds(),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger),
			cron.Recover(cron.DefaultLogger)))
	log.Println("开启定时任务")
	C.AddFunc("@daily", CleanLogTask)
	C.Start()
}
