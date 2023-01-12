package models

import (
	"time"
)

type Permission struct {
	ID        uint      `json:"id"`
	Name      string    `gorm:"size:256;unique" json:"name" filter:"name" validate:"required"`
	URL       string    `gorm:"size:128;index:idx_route,unique" json:"url" filter:"url" validate:"required"`
	Method    string    `gorm:"size:32;index:idx_route,unique" json:"method" filter:"method" validate:"required"` // GET POST PUT DELETE
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Permissions []Permission
