package models

import (
	"database/sql"
	"time"
)

type Log struct {
	ID        uint
	Username  string         `gorm:"size:128"`
	IP        string         `gorm:"size:128"`
	Method    string         `gorm:"size:16"`
	URL       string         `gorm:"size:128"`
	Data      sql.NullString `gorm:"type:text"`
	CreatedAt time.Time
}

type Logs []Log
