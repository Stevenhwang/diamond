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
	respMsg(c, 0, "更新成功！")
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
func ResAssignPerm(c *gin.Context) {
	type resAssign struct {
		Type      string `json:"type" binding:"required"` // 资源类型(Users, Groups, Menus, Permissions)
		Resources []int  `json:"resources"`               // 资源id列表
	}
	assign := &resAssign{}
	if err := c.ShouldBindJSON(assign); err != nil {
		respMsg(c, 1, err.Error())
		return
	}
	role := &models.Role{}
	if result := models.DB.Find(role, c.Param("id")); result.Error != nil {
		respMsg(c, 2, result.Error.Error())
		return
	}
	// 清空原有关联
	models.DB.Model(role).Association(assign.Type).Clear()
	// 添加关联
	switch assign.Type {
	case "Users":
		users := &models.Users{}
		models.DB.Find(users, assign.Resources)
		models.DB.Model(role).Association(assign.Type).Append(users)
	case "Groups":
		groups := &models.Groups{}
		models.DB.Find(groups, assign.Resources)
		models.DB.Model(role).Association(assign.Type).Append(groups)
	case "Menus":
		menus := &models.Menus{}
		models.DB.Find(menus, assign.Resources)
		models.DB.Model(role).Association(assign.Type).Append(menus)
	case "Permissions":
		permissions := &models.Permissions{}
		models.DB.Find(permissions, assign.Resources)
		models.DB.Model(role).Association(assign.Type).Append(permissions)
	}
	respMsg(c, 0, "分配成功！")
}
