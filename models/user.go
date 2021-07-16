package models

import (
	"database/sql"
	"time"

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

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	return
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	return
}

// 在同一个事务中更新数据
func (u *User) AfterUpdate(tx *gorm.DB) (err error) {

	return
}

// 在同一个事务中更新数据
func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	return
}
