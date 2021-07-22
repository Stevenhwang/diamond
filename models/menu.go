package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/nulls"
)

type Menu struct {
	ID        uint         `json:"id"`
	Name      string       `gorm:"size:128;unique" json:"name" filter:"name" binding:"required"`
	Remark    nulls.String `gorm:"size:256" json:"remark" filter:"remark"`
	IsActive  bool         `gorm:"default:true" json:"is_active"`
	Roles     []*Role      `gorm:"many2many:role_menus" json:"-"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type Menus []Menu

func GetMenuList(c *gin.Context) (menus Menus, total int64, err error) {
	menus = Menus{}
	// 使用role_id查找的时候不用分页，也不用filter
	if roleID := c.Query("role_id"); len(roleID) > 0 {
		role := &Role{}
		DB.First(role, roleID)
		total = DB.Model(role).Association("Menus").Count()
		err := DB.Model(role).Association("Menus").Find(menus)
		return menus, total, err
	}
	DB.Model(&Menu{}).Scopes(Filter(Menu{}, c)).Count(&total)
	result := DB.Scopes(Filter(Menu{}, c), Paginate(c)).Find(&menus)
	return menus, total, result.Error
}
