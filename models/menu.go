package models

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Menu struct {
	ID        uint           `json:"id"`
	Name      string         `gorm:"size:128;unique" json:"name"`
	Remark    sql.NullString `gorm:"size:128" json:"remark"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	Roles     []*Role        `gorm:"many2many:role_menus"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type Menus []Menu

func GetMenuList(c *fiber.Ctx) (menus Menus, total int64, err error) {
	menus = Menus{}
	DB.Model(&Menu{}).Scopes(Filter(Menu{}, c)).Count(&total)
	result := DB.Scopes(Filter(Menu{}, c), Paginate(c)).Find(&menus)
	return menus, total, result.Error
}
