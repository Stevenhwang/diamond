package actions

import (
	"diamond/models"
	"diamond/policy"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func getServers(c echo.Context) error {
	var total int64
	servers := models.Servers{}
	is_superuser := c.Get("is_superuser").(bool)
	basequery := models.DB.Model(&models.Server{}).Scopes(models.Filter(models.Server{}, c))
	if !is_superuser {
		uid := c.Get("uid").(string)
		sub := fmt.Sprintf("user::%s", uid)
		var sids models.Servers
		models.DB.Model(&models.Server{}).Select("id", "is_active").Find(&sids)
		var requests [][]interface{}
		for _, s := range sids {
			obj := fmt.Sprintf("server::%d", s.ID)
			requests = append(requests, []interface{}{sub, obj, "server"})
		}
		results, err := policy.Enforcer.BatchEnforce(requests)
		if err != nil {
			return c.JSON(http.StatusOK, H{"code": 1, "message": err.Error()})
		}
		var res []uint
		for i, r := range sids {
			if results[i] && r.IsActive {
				res = append(res, r.ID)
			}
		}
		basequery = basequery.Where("id IN ?", res)
	}
	basequery.Count(&total)
	result := basequery.Scopes(models.Paginate(c)).Find(&servers)
	if result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, H{"code": 0, "data": servers, "total": total})
}

func createServer(c echo.Context) error {
	server := &models.Server{}
	if err := c.Bind(server); err != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": err.Error()})
	}
	if err := c.Validate(server); err != nil {
		return c.JSON(http.StatusOK, H{"code": 2, "message": err.Error()})
	}
	if result := models.DB.Create(server); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 3, "message": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, H{"code": 0, "message": "create server success"})
}

func updateServer(c echo.Context) error {
	server := &models.Server{}
	if result := models.DB.Find(server, c.Param("id")); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": result.Error.Error()})
	}
	if err := c.Bind(server); err != nil {
		return c.JSON(http.StatusOK, H{"code": 2, "message": err.Error()})
	}
	if err := c.Validate(server); err != nil {
		return c.JSON(http.StatusOK, H{"code": 3, "message": err.Error()})
	}
	if result := models.DB.Select("*").Updates(server); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 4, "message": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, H{"code": 0, "message": "update server success"})
}

func deleteServer(c echo.Context) error {
	server := &models.Server{}
	if result := models.DB.Delete(server, c.Param("id")); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, H{"code": 0, "message": "delete server success"})
}
