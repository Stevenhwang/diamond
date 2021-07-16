package models

type Permission struct {
	ID       int
	Name     string `gorm:"size:128"`
	Remark   string `gorm:"size:256"`
	IsActive bool   `gorm:"default:true"`
}
