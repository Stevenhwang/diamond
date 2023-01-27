package models

import "time"

type Cron struct {
	ID        uint      `json:"id"`
	EntryID   int       `json:"entryID"`                                                    // 定时任务EntryID(自动生成)
	Name      string    `gorm:"size:128" json:"name" validate:"required" filter:"name"`     // 定时任务名称
	Target    string    `gorm:"size:128" json:"target" validate:"required" filter:"target"` // ansible目标(服务器ip或服务器分组)
	ScriptID  uint      `json:"script_id" validate:"required"`                              // 关联执行的脚本
	Args      string    `gorm:"size:256" json:"args"`                                       // 执行script时传入的参数(空格分开)
	Spec      string    `gorm:"size:128" json:"spec" validate:"required"`                   // 定时任务时间
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Crons []Cron
