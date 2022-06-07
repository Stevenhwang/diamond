package models

import (
	"errors"
	"time"

	"github.com/gobuffalo/nulls"
	"gorm.io/gorm"
)

type Server struct {
	ID        uint         `json:"id"`
	IP        string       `gorm:"size:128" json:"ip" filter:"ip" validate:"required,ipv4"`
	Hostname  nulls.String `gorm:"size:128" json:"hostname" filter:"ip"`
	Remark    nulls.String `gorm:"size:256" json:"remark" filter:"remark"`
	Port      uint         `gorm:"default:22" json:"port" validate:"required"`
	User      string       `gorm:"size:128" json:"user" validate:"required"`
	AuthType  uint         `json:"auth_type" validate:"required"` // 认证方式：1.password/2.key
	Password  nulls.String `gorm:"size:128" json:"password"`
	KeyID     nulls.UInt32 `json:"key_id"`
	IsActive  bool         `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type Servers []Server

func (s *Server) BeforeCreate(tx *gorm.DB) (err error) {
	if s.AuthType != 1 && s.AuthType != 2 {
		return errors.New("auth type is only 1 or 2")
	}
	if s.AuthType == 1 {
		if len(s.Password.String) == 0 {
			return errors.New("need password")
		}
	} else {
		if !s.KeyID.Valid {
			return errors.New("need key")
		}
	}
	return nil
}

func (s *Server) BeforeUpdate(tx *gorm.DB) (err error) {
	if s.AuthType != 1 && s.AuthType != 2 {
		return errors.New("auth type is only 1 or 2")
	}
	if s.AuthType == 1 {
		if len(s.Password.String) == 0 {
			return errors.New("need password")
		}
	} else {
		if !s.KeyID.Valid {
			return errors.New("need key")
		}
	}
	return nil
}
