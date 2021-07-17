package models

import (
	"time"
)

type Group struct {
	ID        uint
	Name      string   `gorm:"size:128;unique"`
	IsActive  bool     `gorm:"default:true"`
	Roles     []*Role  `gorm:"many2many:role_groups"`
	Servers   []Server `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Groups []Group
