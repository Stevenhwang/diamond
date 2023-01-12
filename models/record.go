package models

import (
	"diamond/misc"
	"os"
	"time"

	"gorm.io/gorm"
)

type Record struct {
	ID        uint      `json:"id"`
	User      string    `gorm:"size:128" json:"user"` // 用户
	IP        string    `gorm:"size:128" json:"ip"`   // 服务器IP
	File      string    `gorm:"size:128" json:"file"` // 记录文件名
	CreatedAt time.Time `json:"created_at"`
}

type Records []Record

func (r *Record) AfterDelete(tx *gorm.DB) (err error) {
	if err := os.Remove(r.File); err != nil {
		misc.Logger.Info().Str("from", "db").Msg("file not exists")
	}
	misc.Logger.Info().Str("from", "db").Msg("file removed")
	return
}
