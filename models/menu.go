package models

import (
	"database/sql"
	"time"
)

type Menu struct {
	ID        uint
	Name      string         `gorm:"size:128;unique"`
	Remark    sql.NullString `gorm:"size:128"`
	IsActive  bool           `gorm:"default:true"`
	Roles     []*Role        `gorm:"many2many:role_menus"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Menus []Menu
