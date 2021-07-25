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
	// 处理password和key更新
	excludeColumns := []string{}
	if len(server.Password.String) == 0 {
		excludeColumns = append(excludeColumns, "password")
	}
	if len(server.Key.String) == 0 {
		excludeColumns = append(excludeColumns, "key")
	}
	if result := models.DB.Select("*").Omit(excludeColumns...).Updates(server); result.Error != nil {
		respMsg(c, 3, result.Error.Error())
		return
	}
	respMsg(c, 0, "更新成功！")
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
