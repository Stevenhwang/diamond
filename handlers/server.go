package handlers

import (
	"diamond/models"

	"github.com/gin-gonic/gin"
)

// 获取服务器列表
func ServerListPerm(c *gin.Context) {
	servers, total, err := models.GetServerList(c)
	if err != nil {
		respMsg(c, 1, err.Error())
		return
	}
	respData(c, 0, servers, total)
}

// 更新服务器信息
func UpdateServerPerm(c *gin.Context) {
	server := &models.Server{}
	if result := models.DB.Find(server, c.Param("id")); result.Error != nil {
		respMsg(c, 1, result.Error.Error())
		return
	}
	if err := c.ShouldBindJSON(server); err != nil {
		respMsg(c, 2, err.Error())
		return
	}
	if result := models.DB.Select("*").Updates(server); result.Error != nil {
		respMsg(c, 3, result.Error.Error())
		return
	}
}

// 新建服务器
func CreateServerPerm(c *gin.Context) {
	server := &models.Server{}
	if err := c.ShouldBindJSON(server); err != nil {
		respMsg(c, 1, err.Error())
		return
	}
	if result := models.DB.Create(server); result.Error != nil {
		respMsg(c, 2, result.Error.Error())
		return
	}
	respMsg(c, 0, "创建成功！")
}

// 删除服务器
func DeleteServerPerm(c *gin.Context) {
	server := &models.Server{}
	if result := models.DB.Delete(server, c.Param("id")); result.Error != nil {
		respMsg(c, 1, result.Error.Error())
		return
	}
	respMsg(c, 0, "删除成功！")
}
