package models

import (
	"errors"
	"time"

	"github.com/Stevenhwang/gommon/nulls"
	"github.com/Stevenhwang/gommon/tools"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID            uint           `json:"id"`
	Username      string         `gorm:"size:128;unique" json:"username" validate:"required"`
	Password      string         `gorm:"size:128" json:"password"`
	Publickey     string         `gorm:"type:text" json:"publickey"` // 公钥，用于免密登录ssh服务器
	Menus         datatypes.JSON `gorm:"type:json" json:"menus"`     // 给用户分配的菜单
	IsActive      bool           `json:"is_active"`                  // 账号是否激活
	LastLoginIP   nulls.String   `gorm:"size:128" json:"last_login_ip"`
	LastLoginTime nulls.Time     `json:"last_login_time"`
	Servers       Servers        `gorm:"many2many:user_servers;"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type Users []User

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if len(u.Password) == 0 {
		return errors.New("password can not be empty")
	}
	pass, err := tools.GeneratePassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = pass
	// generate otp key
	return nil
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if len(u.Password) > 0 {
		pass, err := tools.GeneratePassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = pass
	}
	return nil
}

func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
	return
}

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	return
}
