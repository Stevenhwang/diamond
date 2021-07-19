package models

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Group struct {
	ID        uint      `json:"id"`
	Name      string    `gorm:"size:128;unique" json:"name"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	Roles     []*Role   `gorm:"many2many:role_groups"`
	Servers   []Server  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Groups []Group

func GetGroupList(c *gin.Context) (groups Groups, total int64, err error) {
	groups = Groups{}
	DB.Model(&Group{}).Scopes(Filter(Group{}, c)).Count(&total)
	result := DB.Scopes(Filter(Group{}, c), Paginate(c)).Find(&groups)
	return groups, total, result.Error
}
