package models

import (
	"diamond/utils"
	"strconv"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

type User struct {
	ID            uint         `json:"id"`
	Username      string       `gorm:"size:128;unique" filter:"username" json:"username" validate:"required"`
	Password      string       `gorm:"size:128" json:"password"`
	Email         nulls.String `gorm:"size:128" filter:"email" json:"email" validate:"omitempty,email"`
	Telephone     nulls.String `gorm:"size:20" filter:"telephone" json:"telephone"`
	Department    nulls.String `gorm:"size:128" filter:"department" json:"department"`
	GoogleKey     nulls.String `gorm:"size:256" json:"google_key"`
	IsActive      bool         `gorm:"default:true" json:"is_active"`
	IsSuperuser   bool         `gorm:"default:false" json:"is_superuser"`
	LastLoginIP   nulls.String `gorm:"size:128" filter:"last_login_ip" json:"last_login_ip"`
	LastLoginTime nulls.Time   `json:"last_login_time"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

type Users []User

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if len(u.Password) == 0 {
		u.Password = "12345678" // 默认密码
	}
	pass, err := utils.GeneratePassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = pass
	// generate otp key
	if u.GoogleKey.String == "seed" {
		u.GoogleKey = nulls.NewString("")
	} else {
		var accountName string
		if len(u.Email.String) > 0 {
			accountName = u.Email.String
		} else {
			accountName = u.Username + "@diamond"
		}
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "diamond",
			AccountName: accountName,
		})
		if err != nil {
			return err
		}
		u.GoogleKey = nulls.NewString(key.Secret())
	}
	return nil
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if len(u.Password) > 0 {
		pass, err := utils.GeneratePassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = pass
	}
	// 处理otpkey更新
	otpKey, _ := strconv.Atoi(u.GoogleKey.String)
	if otpKey == 2 { // 重置key
		var accountName string
		if len(u.Email.String) > 0 {
			accountName = u.Email.String
		} else {
			accountName = u.Username + "@diamond"
		}
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "diamond",
			AccountName: accountName,
		})
		if err != nil {
			return err
		}
		u.GoogleKey = nulls.NewString(key.Secret())
	} else if otpKey == 3 { // 清空key
		u.GoogleKey = nulls.NewString("")
	}
	// 处理is_active更新
	if !u.IsActive {
		utils.DelToken(u.ID)
	}
	return nil
}

func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
	return
}

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	return
}
