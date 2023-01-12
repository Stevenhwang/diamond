package models

import (
	"errors"
	"time"

	"github.com/Stevenhwang/gommon/tools"
	"gorm.io/gorm"
)

type Credential struct {
	ID          uint      `json:"id"`
	Name        string    `gorm:"size:128" json:"name" validate:"required"`      // 名称
	AuthType    uint      `json:"auth_type" validate:"required"`                 // 认证类型 1.密码 2.密钥
	AuthUser    string    `gorm:"size:128" json:"auth_user" validate:"required"` // 认证用户
	AuthContent string    `gorm:"type:text" json:"auth_content"`                 // 认证内容，密码或者私钥(没必要用加密的私钥)
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Credentials []Credential

func (s *Credential) BeforeCreate(tx *gorm.DB) (err error) {
	if len(s.AuthContent) > 0 {
		pass := tools.AesEncrypt(s.AuthContent, "0123456789012345")
		s.AuthContent = pass
	}
	return
}

func (s *Credential) BeforeUpdate(tx *gorm.DB) (err error) {
	if len(s.AuthContent) > 0 {
		pass := tools.AesEncrypt(s.AuthContent, "0123456789012345")
		s.AuthContent = pass
	}
	return
}

func (u *Credential) BeforeDelete(tx *gorm.DB) (err error) {
	var count int64
	DB.Model(&Server{}).Where("credential_id = ?", u.ID).Count(&count)
	if count > 0 {
		return errors.New("can not delete because an associated server exists")
	}
	return
}
