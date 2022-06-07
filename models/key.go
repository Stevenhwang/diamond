package models

import (
	"time"

	"github.com/gobuffalo/nulls"
)

type Key struct {
	ID        uint32       `json:"id"`
	Name      string       `gorm:"size:128;unique" json:"name" filter:"name" validate:"required"`
	Remark    nulls.String `gorm:"size:256" json:"remark" filter:"remark"`
	Content   string       `gorm:"type:text" json:"key" validate:"required"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`

	Servers Servers `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Keys []Key
