package models

import "time"

type Menu struct {
	ID        int
	Name      string  `gorm:"size:128"`
	Remark    string  `gorm:"size:128"`
	IsActive  bool    `gorm:"default:true"`
	Roles     []*Role `gorm:"many2many:role_menus"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
