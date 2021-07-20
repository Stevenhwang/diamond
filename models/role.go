package models

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Role struct {
	ID          uint          `json:"id"`
	Name        string        `gorm:"size:128;unique" json:"name" filter:"name" binding:"required"`
	IsActive    bool          `gorm:"default:true" json:"is_active"`
	Users       []*User       `gorm:"many2many:user_roles" json:"-"`
	Permissions []*Permission `gorm:"many2many:role_permissions" json:"-"`
	Menus       []*Menu       `gorm:"many2many:role_menus" json:"-"`
	Groups      []*Group      `gorm:"many2many:role_groups" json:"-"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type Roles []Role

func GetRoleList(c *gin.Context) (roles Roles, total int64, err error) {
	roles = Roles{}
	DB.Model(&Role{}).Scopes(Filter(Role{}, c)).Count(&total)
	result := DB.Scopes(Filter(Role{}, c), Paginate(c)).Find(&roles)
	return roles, total, result.Error
}
