package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint
	UserName      string `gorm:"size:128"`
	PassWord      string `gorm:"size:128"`
	Email         string `gorm:"size:128"`
	Telephone     string `gorm:"size:20"`
	Department    string `gorm:"size:128"`
	GoogleKey     string `gorm:"size:256"`
	IsActive      bool   `gorm:"default:true"`
	IsSuperuser   bool   `gorm:"default:false"`
	LastLoginIP   string `gorm:"size:128"`
	LastLoginTime time.Time
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
