package models

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Group struct {
	ID        uint      `json:"id"`
	Name      string    `gorm:"size:128;unique" json:"name" filter:"name" binding:"required"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	Roles     []*Role   `gorm:"many2many:role_groups" json:"-"`
	Servers   []Server  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Groups []Group

func GetGroupList(c *gin.Context) (groups Groups, total int64, err error) {
	groups = Groups{}
	// 使用role_id查找的时候不用分页，也不用filter
	if roleID := c.Query("role_id"); len(roleID) > 0 {
		role := &Role{}
		DB.First(role, roleID)
		total = DB.Model(role).Association("Groups").Count()
		err := DB.Model(role).Association("Groups").Find(groups)
		return groups, total, err
	}
	DB.Model(&Group{}).Scopes(Filter(Group{}, c)).Count(&total)
	result := DB.Scopes(Filter(Group{}, c), Paginate(c)).Find(&groups)
	return groups, total, result.Error
}
