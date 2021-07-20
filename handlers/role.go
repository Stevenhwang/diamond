package handlers

import (
	"diamond/models"

	"github.com/gin-gonic/gin"
)

// 获取角色列表
func RoleListPerm(c *gin.Context) {
	roles, total, err := models.GetRoleList(c)
	if err != nil {
		respMsg(c, 1, err.Error())
		return
	}
	respData(c, 0, roles, total)
}

// 更新角色信息
func UpdateRolePerm(c *gin.Context) {
	role := &models.Role{}
	if result := models.DB.Find(role, c.Param("id")); result.Error != nil {
		respMsg(c, 1, result.Error.Error())
		return
	}
	if err := c.ShouldBindJSON(role); err != nil {
		respMsg(c, 2, err.Error())
		return
	}
	if result := models.DB.Select("*").Updates(role); result.Error != nil {
		respMsg(c, 3, result.Error.Error())
		return
	}
}

// 新建角色
func CreateRolePerm(c *gin.Context) {
	role := &models.Role{}
	if err := c.ShouldBindJSON(role); err != nil {
		respMsg(c, 1, err.Error())
		return
	}
	if result := models.DB.Create(role); result.Error != nil {
		respMsg(c, 2, result.Error.Error())
		return
	}
	respMsg(c, 0, "创建成功！")
}

// 删除角色
func DeleteRolePerm(c *gin.Context) {
	role := &models.Role{}
	if result := models.DB.Delete(role, c.Param("id")); result.Error != nil {
		respMsg(c, 1, result.Error.Error())
		return
	}
	respMsg(c, 0, "删除成功！")
}

// 给角色分配资源
func ResAssignPerm(c *gin.Context) {}
