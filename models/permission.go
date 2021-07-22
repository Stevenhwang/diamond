package models

import (
	"github.com/gin-gonic/gin"
)

type Permission struct {
	ID       uint    `json:"id"`
	Name     string  `gorm:"size:128;unique" json:"name" filter:"name"`
	Remark   string  `gorm:"size:512" json:"remark" filter:"remark"`
	IsActive bool    `gorm:"default:true" json:"is_active"`
	Roles    []*Role `gorm:"many2many:role_permissions" json:"-"`
}

type Permissions []Permission

func GetPermissionList(c *gin.Context) (permissions Permissions, total int64, err error) {
	permissions = Permissions{}
	// 使用role_id查找的时候不用分页，也不用filter
	if roleID := c.Query("role_id"); len(roleID) > 0 {
		role := &Role{}
		DB.First(role, roleID)
		total = DB.Model(role).Association("Permissions").Count()
		err := DB.Model(role).Association("Permissions").Find(permissions)
		return permissions, total, err
	}
	DB.Model(&Permission{}).Scopes(Filter(Permission{}, c)).Count(&total)
	result := DB.Scopes(Filter(Permission{}, c), Paginate(c)).Find(&permissions)
	return permissions, total, result.Error
}
