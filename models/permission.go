package models

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type Permission struct {
	ID       uint           `json:"id"`
	Name     string         `gorm:"size:128;unique" json:"name"`
	Remark   sql.NullString `gorm:"size:256" json:"remark"`
	IsActive bool           `gorm:"default:true" json:"is_active"`
	Roles    []*Role        `gorm:"many2many:role_permissions"`
}

type Permissions []Permission

func GetPermissionList(c *fiber.Ctx) (permissions Permissions, total int64, err error) {
	permissions = Permissions{}
	DB.Scopes(Filter(Permission{}, c)).Count(&total)
	result := DB.Scopes(Filter(Permission{}, c), Paginate(c)).Find(&permissions)
	return permissions, total, result.Error
}
