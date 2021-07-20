package handlers

import (
	"diamond/models"

	"github.com/gin-gonic/gin"
)

// 获取前端菜单列表
func MenuListPerm(c *gin.Context) {
	menus, total, err := models.GetMenuList(c)
	if err != nil {
		respMsg(c, 1, err.Error())
		return
	}
	respData(c, 0, menus, total)
}

// 更新前端菜单信息
func UpdateMenuPerm(c *gin.Context) {
	menu := &models.Menu{}
	if result := models.DB.Find(menu, c.Param("id")); result.Error != nil {
		respMsg(c, 1, result.Error.Error())
		return
	}
	if err := c.ShouldBindJSON(menu); err != nil {
		respMsg(c, 2, err.Error())
		return
	}
	if result := models.DB.Select("*").Updates(menu); result.Error != nil {
		respMsg(c, 3, result.Error.Error())
		return
	}
}

// 新建前端菜单
func CreateMenuPerm(c *gin.Context) {
	menu := &models.Menu{}
	if err := c.ShouldBindJSON(menu); err != nil {
		respMsg(c, 1, err.Error())
		return
	}
	if result := models.DB.Create(menu); result.Error != nil {
		respMsg(c, 2, result.Error.Error())
		return
	}
	respMsg(c, 0, "创建成功！")
}

// 删除前端菜单
func DeleteMenuPerm(c *gin.Context) {
	menu := &models.Menu{}
	if result := models.DB.Delete(menu, c.Param("id")); result.Error != nil {
		respMsg(c, 1, result.Error.Error())
		return
	}
	respMsg(c, 0, "删除成功！")
}
