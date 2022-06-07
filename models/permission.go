package models

// import (
// 	"time"

// 	"github.com/gobuffalo/nulls"
// )

// type Permission struct {
// 	ID        uint         `json:"id"`
// 	Name      string       `gorm:"size:128;unique" json:"name" filter:"name" validate:"required"`
// 	Typ       string       `gorm:"size:64" json:"typ" filter:"typ" validate:"required"` // route/menu/server
// 	Object    string       `gorm:"size:128" json:"object" validate:"required"`          // "GET /api/users"/system/game_servers
// 	Remark    nulls.String `gorm:"size:256" json:"remark" filter:"remark"`
// 	IsActive  bool         `gorm:"default:true" json:"is_active"`
// 	CreatedAt time.Time    `json:"created_at"`
// 	UpdatedAt time.Time    `json:"updated_at"`
// }

// type Permissions []Permission
