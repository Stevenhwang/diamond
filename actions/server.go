package actions

import (
	"diamond/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// 获取servers, 根据用户分配的服务器获取
func getServers(c echo.Context) error {
	var total int64
	servers := models.Servers{}
	baseQuery := models.DB.Order("created_at desc")
	username := c.Get("username")
	if username.(string) != "admin" {
		user := models.User{}
		uid := c.Get("uid").(uint)
		models.DB.Preload("Servers", func(db *gorm.DB) *gorm.DB {
			return db.Where(db.Order("created_at desc")).Where(db.Scopes(models.AnyFilter(models.Server{}, c)))
		}).First(&user, uid)
		total = int64(len(user.Servers))
		var sids []uint
		for _, s := range user.Servers {
			sids = append(sids, s.ID)
		}
		if res := baseQuery.Scopes(models.Paginate(c)).Where("id IN ?", sids).Find(&servers); res.Error != nil {
			return echo.NewHTTPError(400, res.Error.Error())
		}
	} else {
		if res := baseQuery.Scopes(models.AnyFilter(models.Server{}, c)).Model(&models.Server{}).Count(&total); res.Error != nil {
			return echo.NewHTTPError(400, res.Error.Error())
		}
		if res := baseQuery.Scopes(models.Paginate(c), models.AnyFilter(models.Server{}, c)).Find(&servers); res.Error != nil {
			return echo.NewHTTPError(400, res.Error.Error())
		}
	}
	return c.JSON(200, echo.Map{"success": true, "data": servers, "total": total})
}

func createServer(c echo.Context) error {
	server := models.Server{}
	if err := c.Bind(&server); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if err := c.Validate(&server); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if result := models.DB.Create(&server); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true})
}

func updateServer(c echo.Context) error {
	server := models.Server{}
	if result := models.DB.First(&server, c.Param("id")); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	if err := c.Bind(&server); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if err := c.Validate(&server); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if result := models.DB.Save(&server); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true})
}

func deleteServer(c echo.Context) error {
	server := models.Server{}
	if res := models.DB.Delete(&server, c.Param("id")); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true})
}

func getRecords(c echo.Context) error {
	var total int64
	records := models.Records{}
	if res := models.DB.Model(&models.Record{}).Count(&total); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	if res := models.DB.Scopes(models.Paginate(c)).Order("created_at desc").Find(&records); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true, "data": records, "total": total})
}
