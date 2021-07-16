package models

import (
	"database/sql"
	"time"
)

type Server struct {
	ID        uint
	IP        string         `gorm:"size:128"`
	Remark    sql.NullString `gorm:"size:128"`
	Port      int
	User      string         `gorm:"size:128"`
	AuthType  int            // 1密码验证 2密钥验证
	Password  sql.NullString `gorm:"size:128"`
	Key       sql.NullString `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
