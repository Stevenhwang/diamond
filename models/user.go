package models

import (
	"time"
)

type User struct {
	ID            uint
	UserName      string `gorm:"size:128"`
	PassWord      string `gorm:"size:128"`
	Email         string `gorm:"size:128"`
	Telephone     string `gorm:"size:20"`
	Department    string `gorm:"size:128"`
	GoogleKey     string `gorm:"size:128"`
	IsActive      bool
	IsSuperuser   bool
	LastLoginIP   string `gorm:"size:128"`
	LastLoginTime time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
