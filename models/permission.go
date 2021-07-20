package models

import (
	"github.com/gin-gonic/gin"
)

type Permission struct {
	ID       uint    `json:"id"`
	Name     string  `gorm:"size:128;unique" json:"name"`
	Remark   string  `gorm:"size:256" json:"remark"`
	IsActive bool    `gorm:"default:true" json:"is_active"`
	Roles    []*Role `gorm:"many2many:role_permissions" json:"-"`
}

type Permissions []Permission

func GetPermissionList(c *gin.Context) (permissions Permissions, total int64, err error) {
	permissions = Permissions{}
	DB.Model(&Permission{}).Scopes(Filter(Permission{}, c)).Count(&total)
	result := DB.Scopes(Filter(Permission{}, c), Paginate(c)).Find(&permissions)
	return permissions, total, result.Error
}
