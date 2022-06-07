package actions

import (
	"diamond/models"
	"diamond/policy"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// FIXME 数据权限过滤
func getGroups(c echo.Context) error {
	var total int64
	groups := models.Groups{}
	models.DB.Model(&models.Group{}).Scopes(models.Filter(models.Group{}, c)).Count(&total)
	result := models.DB.Scopes(models.Filter(models.Group{}, c), models.Paginate(c)).Find(&groups)
	if result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, H{"code": 0, "data": groups, "total": total})
}

func createGroup(c echo.Context) error {
	group := &models.Group{}
	if err := c.Bind(group); err != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": err.Error()})
	}
	if err := c.Validate(group); err != nil {
		return c.JSON(http.StatusOK, H{"code": 2, "message": err.Error()})
	}
	if result := models.DB.Create(group); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 3, "message": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, H{"code": 0, "message": "create group success"})
}

func updateGroup(c echo.Context) error {
	group := &models.Group{}
	if result := models.DB.Find(group, c.Param("id")); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": result.Error.Error()})
	}
	if err := c.Bind(group); err != nil {
		return c.JSON(http.StatusOK, H{"code": 2, "message": err.Error()})
	}
	if err := c.Validate(group); err != nil {
		return c.JSON(http.StatusOK, H{"code": 3, "message": err.Error()})
	}
	if result := models.DB.Select("*").Updates(group); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 4, "message": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, H{"code": 0, "message": "update group success"})
}

func deleteGroup(c echo.Context) error {
	group := &models.Group{}
	if result := models.DB.Delete(group, c.Param("id")); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": result.Error.Error()})
	}
	// delete policy for role
	obj := fmt.Sprintf("group::%d", group.ID)
	policy.Enforcer.DeleteRole(obj)
	return c.JSON(http.StatusOK, H{"code": 0, "message": "delete group success"})
}
