package models

import (
	"time"

	"github.com/gobuffalo/nulls"
)

type Group struct {
	ID        uint         `json:"id"`
	Name      string       `gorm:"size:128;unique" json:"name" filter:"name" validate:"required"`
	Remark    nulls.String `gorm:"size:256" json:"remark" filter:"remark"`
	IsActive  bool         `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type Groups []Group
