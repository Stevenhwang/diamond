package actions

import (
	"diamond/models"
	"diamond/policy"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func getGroups(c echo.Context) error {
	var total int64
	groups := models.Groups{}
	is_superuser := c.Get("is_superuser").(bool)
	basequery := models.DB.Model(&models.Group{}).Scopes(models.Filter(models.Group{}, c))
	if !is_superuser {
		uid := c.Get("uid").(string)
		sub := fmt.Sprintf("user::%s", uid)
		var gids models.Groups
		models.DB.Model(&models.Group{}).Select("id", "is_active").Find(&gids)
		var requests [][]interface{}
		for _, s := range gids {
			obj := fmt.Sprintf("group::%d", s.ID)
			requests = append(requests, []interface{}{sub, obj, "group"})
		}
		results, err := policy.Enforcer.BatchEnforce(requests)
		if err != nil {
			return c.JSON(http.StatusOK, H{"code": 1, "message": err.Error()})
		}
		var res []uint
		for i, r := range gids {
			if results[i] && r.IsActive {
				res = append(res, r.ID)
			}
		}
		basequery = basequery.Where("id IN ?", res)
	}
	basequery.Count(&total)
	result := basequery.Scopes(models.Paginate(c)).Find(&groups)
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
