package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/nulls"
)

type Server struct {
	ID        uint         `json:"id"`
	IP        string       `gorm:"size:128" json:"ip" filter:"ip" binding:"required,ipv4"`
	Hostname  nulls.String `gorm:"size:128" json:"hostname" filter:"ip"`
	Remark    nulls.String `gorm:"size:256" json:"remark" filter:"remark"`
	Port      int          `gorm:"default:22" json:"port" binding:"required"`
	User      string       `gorm:"size:128" json:"user" binding:"required"`
	Password  nulls.String `gorm:"size:128" json:"password"`
	Key       nulls.String `gorm:"type:text" json:"key"`
	GroupID   nulls.Int    `json:"group_id"`
	Records   Records      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	IsActive  bool         `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type Servers []Server

func GetServerList(c *gin.Context) (servers Servers, total int64, err error) {
	servers = Servers{}
	uid := c.GetUint("user_id")
	isSuperuser := c.GetBool("is_superuser")
	if isSuperuser {
		baseQuery := DB.Model(&Server{}).Scopes(Filter(Server{}, c))
		if groupID := c.Query("group_id"); len(groupID) > 0 {
			baseQuery.Where("group_id = ?", groupID).Count(&total)
			result := baseQuery.Scopes(Paginate(c)).Where("group_id = ?", groupID).Find(&servers)
			return servers, total, result.Error
		}
		baseQuery.Count(&total)
		result := baseQuery.Scopes(Paginate(c)).Find(&servers)
		return servers, total, result.Error
	}
	user := &User{}
	serverMap := map[uint]bool{}
	DB.Preload("Roles.Groups.Servers").First(user, uid)
	for _, role := range user.Roles {
		if role.IsActive {
			for _, group := range role.Groups {
				if group.IsActive {
					for _, server := range group.Servers {
						if server.IsActive {
							serverMap[server.ID] = true
						}
					}
				}
			}
		}
	}
	sIDList := []int{}
	for k := range serverMap {
		sIDList = append(sIDList, int(k))
	}
	DB.Model(&Server{}).Scopes(Filter(Server{}, c)).Where("id IN ?", sIDList).Count(&total)
	result := DB.Scopes(Filter(Server{}, c), Paginate(c)).Find(&servers, sIDList)
	return servers, total, result.Error
}
