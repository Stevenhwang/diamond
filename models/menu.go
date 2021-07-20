package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/nulls"
)

type Menu struct {
	ID        uint         `json:"id"`
	Name      string       `gorm:"size:128;unique" json:"name" filter:"name" binding:"required"`
	Remark    nulls.String `gorm:"size:128" json:"remark" filter:"remark"`
	IsActive  bool         `gorm:"default:true" json:"is_active"`
	Roles     []*Role      `gorm:"many2many:role_menus" json:"-"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type Menus []Menu

func GetMenuList(c *gin.Context) (menus Menus, total int64, err error) {
	menus = Menus{}
	DB.Model(&Menu{}).Scopes(Filter(Menu{}, c)).Count(&total)
	result := DB.Scopes(Filter(Menu{}, c), Paginate(c)).Find(&menus)
	return menus, total, result.Error
}
