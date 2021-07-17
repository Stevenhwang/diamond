package models

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type User struct {
	ID            uint
	UserName      string         `gorm:"size:128"`
	PassWord      string         `gorm:"size:128"`
	Email         sql.NullString `gorm:"size:128"`
	Telephone     sql.NullString `gorm:"size:20"`
	Department    sql.NullString `gorm:"size:128"`
	GoogleKey     sql.NullString `gorm:"size:256"`
	IsActive      bool           `gorm:"default:true"`
	IsSuperuser   bool           `gorm:"default:false"`
	LastLoginIP   sql.NullString `gorm:"size:128"`
	LastLoginTime sql.NullTime
	Roles         []*Role `gorm:"many2many:user_roles"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Users []User

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	return
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	return
}

func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
	return
}

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	return
}

func (u *User) GetUserList(c *fiber.Ctx) (users Users, total int64, err error) {
	users = Users{}
	DB.Scopes(Filter(User{}, c)).Count(&total)
	result := DB.Scopes(Filter(User{}, c), Paginate(c)).Find(&users)
	return users, total, result.Error
}
