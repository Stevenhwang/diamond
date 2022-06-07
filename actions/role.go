package actions

import (
	"diamond/models"
	"diamond/policy"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func getRoles(c echo.Context) error {
	var total int64
	roles := models.Roles{}
	models.DB.Model(&models.Role{}).Scopes(models.Filter(models.Role{}, c)).Count(&total)
	result := models.DB.Scopes(models.Filter(models.Role{}, c), models.Paginate(c)).Find(&roles)
	if result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, H{"code": 0, "data": roles, "total": total})
}

func createRole(c echo.Context) error {
	role := &models.Role{}
	if err := c.Bind(role); err != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": err.Error()})
	}
	if err := c.Validate(role); err != nil {
		return c.JSON(http.StatusOK, H{"code": 2, "message": err.Error()})
	}
	if result := models.DB.Create(role); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 3, "message": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, H{"code": 0, "message": "create role success"})
}

func updateRole(c echo.Context) error {
	role := &models.Role{}
	if result := models.DB.Find(role, c.Param("id")); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": result.Error.Error()})
	}
	if err := c.Bind(role); err != nil {
		return c.JSON(http.StatusOK, H{"code": 2, "message": err.Error()})
	}
	if err := c.Validate(role); err != nil {
		return c.JSON(http.StatusOK, H{"code": 3, "message": err.Error()})
	}
	if result := models.DB.Select("*").Updates(role); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 4, "message": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, H{"code": 0, "message": "update role success"})
}

func deleteRole(c echo.Context) error {
	role := &models.Role{}
	if result := models.DB.Delete(role, c.Param("id")); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": result.Error.Error()})
	}
	// delete policy for role
	sub := fmt.Sprintf("role::%d", role.ID)
	policy.Enforcer.DeleteRole(sub)
	return c.JSON(http.StatusOK, H{"code": 0, "message": "delete role success"})
}
