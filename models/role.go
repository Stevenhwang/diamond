package models

import "time"

type Role struct {
	ID        int
	Name      string `gorm:"size:128"`
	IsActive  bool   `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
