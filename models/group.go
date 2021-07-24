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
	Servers   Servers   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Groups []Group

func GetGroupList(c *gin.Context) (groups Groups, total int64, err error) {
	groups = Groups{}
	uid := c.GetUint("user_id")
	isSuperuser := c.GetBool("is_superuser")
	if isSuperuser {
		if roleID := c.Query("role_id"); len(roleID) > 0 {
			role := &Role{}
			DB.First(role, roleID)
			total = DB.Model(role).Scopes(Filter(Group{}, c)).Association("Groups").Count()
			err := DB.Model(role).Scopes(Filter(Group{}, c), Paginate(c)).Association("Groups").Find(groups)
			return groups, total, err
		}
		DB.Model(&Group{}).Scopes(Filter(Group{}, c)).Count(&total)
		result := DB.Scopes(Filter(Group{}, c), Paginate(c)).Find(&groups)
		return groups, total, result.Error
	}
	user := &User{}
	groupMap := map[uint]string{}
	DB.Preload("Roles.Groups").First(user, uid)
	for _, role := range user.Roles {
		if role.IsActive {
			for _, group := range role.Groups {
				if group.IsActive {
					groupMap[group.ID] = group.Name
				}
			}
		}
	}
	gIDList := []int{}
	for k := range groupMap {
		gIDList = append(gIDList, int(k))
	}
	DB.Model(&Group{}).Scopes(Filter(Group{}, c)).Where("id IN ?", gIDList).Count(&total)
	result := DB.Scopes(Filter(Group{}, c), Paginate(c)).Find(&groups, gIDList)
	return groups, total, result.Error
}
