package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/nulls"
)

type Server struct {
	ID        uint         `json:"id"`
	IP        string       `gorm:"size:128" json:"ip" filter:"ip" binding:"required,ipv4"`
	Remark    nulls.String `gorm:"size:256" json:"remark" filter:"remark"`
	Port      int          `json:"port" binding:"required"`
	User      string       `gorm:"size:128" json:"user" binding:"required"`
	Password  nulls.String `gorm:"size:128" json:"password"`
	Key       nulls.String `gorm:"type:text" json:"key"`
	GroupID   nulls.Int    `json:"group_id"`
	IsActive  bool         `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type Servers []Server

func GetServerList(c *gin.Context) (servers Servers, total int64, err error) {
	servers = Servers{}
	// 使用group_id查找的时候不用分页，也不用filter
	if groupID := c.Query("group_id"); len(groupID) > 0 {
		DB.Model(&Server{}).Where("group_id = ?", groupID).Count(&total)
		result := DB.Model(&Server{}).Where("group_id = ?", groupID).Find(&servers)
		return servers, total, result.Error
	}
	DB.Model(&Server{}).Scopes(Filter(Server{}, c)).Count(&total)
	result := DB.Scopes(Filter(Server{}, c), Paginate(c)).Find(&servers)
	return servers, total, result.Error
}
