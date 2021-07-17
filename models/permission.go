package models

import "database/sql"

type Permission struct {
	ID       uint
	Name     string         `gorm:"size:128;unique"`
	Remark   sql.NullString `gorm:"size:256"`
	IsActive bool           `gorm:"default:true"`
	Roles    []*Role        `gorm:"many2many:role_permissions"`
}

type Permissions []Permission
