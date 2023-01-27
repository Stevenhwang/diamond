package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Script struct {
	ID        uint      `json:"id"`
	Name      string    `gorm:"size:256;unique" json:"name" validate:"required" filter:"name"` //脚本名称
	Content   string    `gorm:"type:text" json:"content" validate:"required"`                  //脚本内容
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Scripts []Script

func (s *Script) BeforeDelete(tx *gorm.DB) (err error) {
	// 删除之前要检查是否有任务或定时任务在依赖这个脚本
	var countTask int64
	DB.Model(&Task{}).Where("script_id = ?", s.ID).Count(&countTask)
	if countTask > 0 {
		return errors.New("can not delete because an associated task exists")
	}
	var countCron int64
	DB.Model(&Cron{}).Where("script_id = ?", s.ID).Count(&countCron)
	if countCron > 0 {
		return errors.New("can not delete because an associated cron exists")
	}
	return
}
