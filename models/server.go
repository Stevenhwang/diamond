package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/nulls"
)

type Server struct {
	ID        uint         `json:"id"`
	IP        string       `gorm:"size:128" json:"ip" filter:"ip" binding:"required"`
	Remark    nulls.String `gorm:"size:128" json:"remark" filter:"remark"`
	Port      int          `json:"port" binding:"required"`
	User      string       `gorm:"size:128" json:"user" binding:"required"`
	AuthType  int          `json:"auth_type" binding:"required"` // 1密码验证 2密钥验证
	Password  nulls.String `gorm:"size:128" json:"password"`
	Key       nulls.String `gorm:"type:text" json:"key"`
	GroupID   nulls.Int    `json:"group_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type Servers []Server

func GetServerList(c *gin.Context) (servers Servers, total int64, err error) {
	servers = Servers{}
	DB.Model(&Server{}).Scopes(Filter(Server{}, c)).Count(&total)
	result := DB.Scopes(Filter(Server{}, c), Paginate(c)).Find(&servers)
	return servers, total, result.Error
}
