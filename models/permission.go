package models

type Permission struct {
	ID       uint
	Name     string  `gorm:"size:128"`
	Remark   string  `gorm:"size:256"`
	IsActive bool    `gorm:"default:true"`
	Roles    []*Role `gorm:"many2many:role_permissions"`
}
