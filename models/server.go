package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/nulls"
)

type Server struct {
	ID        uint         `json:"id"`
	IP        string       `gorm:"size:128" json:"ip"`
	Remark    nulls.String `gorm:"size:128" json:"remark"`
	Port      int          `json:"port"`
	User      string       `gorm:"size:128" json:"user"`
	AuthType  int          `json:"auth_type"` // 1密码验证 2密钥验证
	Password  nulls.String `gorm:"size:128" json:"password"`
	Key       nulls.String `gorm:"type:text" json:"key"`
	GroupID   uint
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Servers []Server

func GetServerList(c *gin.Context) (servers Servers, total int64, err error) {
	servers = Servers{}
	DB.Model(&Server{}).Scopes(Filter(Server{}, c)).Count(&total)
	result := DB.Scopes(Filter(Server{}, c), Paginate(c)).Find(&servers)
	return servers, total, result.Error
}
