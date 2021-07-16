package models

import (
	"time"
)

type Group struct {
	ID        int
	Name      string  `gorm:"size:128"`
	IsActive  bool    `gorm:"default:true"`
	Roles     []*Role `gorm:"many2many:role_groups"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
