package models

import "time"

type Cron struct {
	ID        uint      `json:"id"`
	EntryID   int       `json:"entryID"`                                                       // 定时任务EntryID(自动生成)
	Name      string    `gorm:"size:128" json:"name" validate:"required" filter:"name"`        // 定时任务名称
	Command   string    `gorm:"type:text" json:"command" validate:"required" filter:"command"` // 执行的命令
	Spec      string    `gorm:"size:128" json:"spec" validate:"required"`                      // 定时任务时间
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Crons []Cron
