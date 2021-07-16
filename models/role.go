package models

import "time"

type Role struct {
	ID          int
	Name        string        `gorm:"size:128"`
	IsActive    bool          `gorm:"default:true"`
	Users       []*User       `gorm:"many2many:user_roles"`
	Permissions []*Permission `gorm:"many2many:role_permissions"`
	Menus       []*Menu       `gorm:"many2many:role_menus"`
	Groups      []*Group      `gorm:"many2many:role_groups"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
