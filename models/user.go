package models

import (
	"database/sql"
	"strconv"
	"time"

	"diamond/utils.go"

	"github.com/gofiber/fiber/v2"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

type User struct {
	ID            uint
	Username      string         `gorm:"size:128" filter:"username"`
	Password      string         `gorm:"size:128"`
	Email         sql.NullString `gorm:"size:128" filter:"email"`
	Telephone     sql.NullString `gorm:"size:20" filter:"telephone"`
	Department    sql.NullString `gorm:"size:128" filter:"department"`
	GoogleKey     sql.NullString `gorm:"size:256"`
	IsActive      bool           `gorm:"default:true"`
	IsSuperuser   bool           `gorm:"default:false"`
	LastLoginIP   sql.NullString `gorm:"size:128" filter:"last_login_ip"`
	LastLoginTime sql.NullTime
	Roles         []*Role `gorm:"many2many:user_roles"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Users []User

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
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
	u.GoogleKey = sql.NullString{String: key.Secret(), Valid: true}
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
		u.GoogleKey = sql.NullString{String: key.Secret(), Valid: true}
	} else if otpKey == 3 {
		u.GoogleKey = sql.NullString{String: "", Valid: true}
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

func GetUserList(c *fiber.Ctx) (users Users, total int64, err error) {
	users = Users{}
	DB.Scopes(Filter(User{}, c)).Count(&total)
	result := DB.Scopes(Filter(User{}, c), Paginate(c)).Find(&users)
	return users, total, result.Error
}
