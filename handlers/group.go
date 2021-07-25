package handlers

import (
	"diamond/models"

	"github.com/gin-gonic/gin"
)

// 获取服务器分组列表
func GroupListPerm(c *gin.Context) {
	groups, total, err := models.GetGroupList(c)
	if err != nil {
		respMsg(c, 1, err.Error())
		return
	}
	respData(c, 0, groups, total)
}

// 更新服务器分组信息
func UpdateGroupPerm(c *gin.Context) {
	group := &models.Group{}
	if result := models.DB.Find(group, c.Param("id")); result.Error != nil {
		respMsg(c, 1, result.Error.Error())
		return
	}
	if err := c.ShouldBindJSON(group); err != nil {
		respMsg(c, 2, err.Error())
		return
	}
	if result := models.DB.Select("*").Updates(group); result.Error != nil {
		respMsg(c, 3, result.Error.Error())
		return
	}
	respMsg(c, 0, "更新成功！")
}

// 新建服务器分组
func CreateGroupPerm(c *gin.Context) {
	group := &models.Group{}
	if err := c.ShouldBindJSON(group); err != nil {
		respMsg(c, 1, err.Error())
		return
	}
	if result := models.DB.Create(group); result.Error != nil {
		respMsg(c, 2, result.Error.Error())
		return
	}
	respMsg(c, 0, "创建成功！")
}

// 删除服务器分组
func DeleteGroupPerm(c *gin.Context) {
	group := &models.Group{}
	if result := models.DB.Delete(group, c.Param("id")); result.Error != nil {
		respMsg(c, 1, result.Error.Error())
		return
	}
	respMsg(c, 0, "删除成功！")
}

// 给服务器组添加服务器
func ServersAssignPerm(c *gin.Context) {
	type serversAssign struct {
		Servers []int `json:"servers"`
	}
	sa := &serversAssign{}
	if err := c.ShouldBindJSON(sa); err != nil {
		respMsg(c, 1, err.Error())
		return
	}
	group := &models.Group{}
	if result := models.DB.Find(group, c.Param("id")); result.Error != nil {
		respMsg(c, 2, result.Error.Error())
		return
	}
	// 清空原有关联
	models.DB.Model(group).Association("Servers").Clear()
	// 添加关联
	servers := &models.Servers{}
	models.DB.Find(servers, sa.Servers)
	models.DB.Model(group).Association("Servers").Append(servers)
	respMsg(c, 0, "分配成功！")
}
