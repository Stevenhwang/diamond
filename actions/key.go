package actions

import (
	"diamond/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func getKeys(c echo.Context) error {
	var total int64
	keys := models.Keys{}
	models.DB.Model(&models.Key{}).Scopes(models.Filter(models.Key{}, c)).Count(&total)
	result := models.DB.Scopes(models.Filter(models.Key{}, c), models.Paginate(c)).Find(&keys)
	if result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, H{"code": 0, "data": keys, "total": total})
}

func createKey(c echo.Context) error {
	key := &models.Key{}
	if err := c.Bind(key); err != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": err.Error()})
	}
	if err := c.Validate(key); err != nil {
		return c.JSON(http.StatusOK, H{"code": 2, "message": err.Error()})
	}
	if result := models.DB.Create(key); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 3, "message": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, H{"code": 0, "message": "create key success"})
}

func updateKey(c echo.Context) error {
	key := &models.Key{}
	if result := models.DB.Find(key, c.Param("id")); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": result.Error.Error()})
	}
	if err := c.Bind(key); err != nil {
		return c.JSON(http.StatusOK, H{"code": 2, "message": err.Error()})
	}
	if err := c.Validate(key); err != nil {
		return c.JSON(http.StatusOK, H{"code": 3, "message": err.Error()})
	}
	if result := models.DB.Select("*").Updates(key); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 4, "message": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, H{"code": 0, "message": "update key success"})
}

func deleteKey(c echo.Context) error {
	key := &models.Key{}
	if result := models.DB.Delete(key, c.Param("id")); result.Error != nil {
		return c.JSON(http.StatusOK, H{"code": 1, "message": result.Error.Error()})
	}
	return c.JSON(http.StatusOK, H{"code": 0, "message": "delete key success"})
}
