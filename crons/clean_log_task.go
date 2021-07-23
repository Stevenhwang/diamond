package crons

import (
	"diamond/models"
	"log"
	"time"
)

// CleanLogTask 定期清理系统日志任务(保留最近15天)
func CleanLogTask() {
	tmpTime := time.Now().AddDate(0, 0, -15)
	log.Printf("开始清理%v之前的日志", tmpTime)
	if result := models.DB.Where("created_at <= ?", tmpTime).Delete(models.Log{}); result.Error != nil {
		log.Println(result.Error)
	}
	log.Println("清理日志完成！")
}
