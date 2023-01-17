package models

import (
	"time"
)

type Task struct {
	ID        uint      `json:"id"`
	Name      string    `gorm:"size:256;unique" json:"name" validate:"required"` // 名称
	Command   string    `gorm:"size:512" json:"command" validate:"required"`     // 执行的命令(脚本单独维护)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Tasks []Task

type TaskHistory struct {
	ID        uint      `json:"id"`
	TaskName  string    `gorm:"size:256" json:"task_name"`
	User      string    `gorm:"size:128" json:"user"`                    // 执行者
	FromIP    string    `gorm:"size:128" json:"from_ip"`                 // from IP
	Success   bool      `json:"success"`                                 // 执行成功、失败
	Content   string    `gorm:"type:text" json:"remark" filter:"remark"` // 执行内容
	CreatedAt time.Time `json:"created_at"`
}

type TaskHistorys []TaskHistory
