package models

import (
	"errors"
	"strconv"
	"time"

	"diamond/utils.go"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/nulls"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

type User struct {
	ID            uint         `json:"id"`
	Username      string       `gorm:"size:128;unique" filter:"username" json:"username"`
	Password      string       `gorm:"size:128" json:"password"`
	Email         nulls.String `gorm:"size:128" filter:"email" json:"email"`
	Telephone     nulls.String `gorm:"size:20" filter:"telephone" json:"telephone"`
	Department    nulls.String `gorm:"size:128" filter:"department" json:"department"`
	GoogleKey     nulls.String `gorm:"size:256" json:"google_key"`
	IsActive      bool         `gorm:"default:true" json:"is_active"`
	IsSuperuser   bool         `gorm:"default:false" json:"is_superuser"`
	LastLoginIP   nulls.String `gorm:"size:128" filter:"last_login_ip" json:"last_login_ip"`
	LastLoginTime nulls.Time   `json:"last_login_time"`
	Roles         []*Role      `gorm:"many2many:user_roles"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

type Users []User

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if len(u.Username) == 0 {
		return errors.New("用户名不能为空")
	}
	// set password
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
			accountName = u.Username + "@example.com"
		}
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      u.Username,
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
	if len(u.Username) == 0 {
		return errors.New("用户名不能为空")
	}
	if len(u.Password) > 0 {
		pass, err := utils.GeneratePassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = pass
	}
	// 处理otpkey更新
	otpKey, _ := strconv.Atoi(u.GoogleKey.String)
	if otpKey == 2 {
		var accountName string
		if len(u.Email.String) > 0 {
			accountName = u.Email.String
		} else {
			accountName = u.Username + "@example.com"
		}
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      u.Username,
			AccountName: accountName,
		})
		if err != nil {
			return err
		}
		u.GoogleKey = nulls.NewString(key.Secret())
	} else if otpKey == 3 {
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

func GetUserList(c *gin.Context) (users Users, total int64, err error) {
	users = Users{}
	DB.Model(&User{}).Scopes(Filter(User{}, c)).Count(&total)
	result := DB.Scopes(Filter(User{}, c), Paginate(c)).Omit("password").Find(&users)
	return users, total, result.Error
}
