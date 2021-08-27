package crons

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

var C *cron.Cron

func Tick() {
	log.Println(time.Now())
}

// 初始化定时任务
func init() {
	C = cron.New(cron.WithSeconds(),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger),
			cron.Recover(cron.DefaultLogger)))
	C.AddFunc("@daily", CleanLogTask)
	// C.AddFunc("@every 3s", Tick)
}
