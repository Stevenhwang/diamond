package handlers

import (
	"diamond/models"

	"github.com/gin-gonic/gin"
)

// 获取权限列表
func PermissionListPerm(c *gin.Context) {
	perms, total, err := models.GetPermissionList(c)
	if err != nil {
		respMsg(c, 1, err.Error())
		return
	}
	respData(c, 0, perms, total)
}

// 更新权限信息
func UpdatePermissionPerm(c *gin.Context) {
	perm := &models.Permission{}
	if result := models.DB.Find(perm, c.Param("id")); result.Error != nil {
		respMsg(c, 1, result.Error.Error())
		return
	}
	if perm.IsActive {
		perm.IsActive = false
		models.DB.Save(perm)
		respMsg(c, 0, "权限禁用成功！")
	} else {
		perm.IsActive = true
		models.DB.Save(perm)
		respMsg(c, 0, "权限启用成功！")
	}
}

// 新建权限
func CreatePermission(c *gin.Context) {}

// 删除权限
func DeletePermission(c *gin.Context) {}
