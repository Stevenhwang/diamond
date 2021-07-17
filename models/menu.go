package models

import (
	"database/sql"
	"time"
)

type Menu struct {
	ID        uint           `json:"id"`
	Name      string         `gorm:"size:128;unique" json:"name"`
	Remark    sql.NullString `gorm:"size:128" json:"remark"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	Roles     []*Role        `gorm:"many2many:role_menus"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type Menus []Menu
