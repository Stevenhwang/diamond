package actions

import (
	"diamond/models"

	"github.com/labstack/echo/v4"
)

func getScripts(c echo.Context) error {
	var total int64
	scripts := models.Scripts{}
	if res := models.DB.Model(&models.Script{}).Count(&total); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	if res := models.DB.Scopes(models.Paginate(c), models.AnyFilter(models.Script{}, c)).Find(&scripts); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true, "data": scripts, "total": total})
}

func createScript(c echo.Context) error {
	script := models.Script{}
	if err := c.Bind(&script); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if err := c.Validate(&script); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if result := models.DB.Create(&script); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true})
}

func updateScript(c echo.Context) error {
	script := models.Script{}
	if result := models.DB.First(&script, c.Param("id")); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	if err := c.Bind(&script); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if err := c.Validate(&script); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if result := models.DB.Save(&script); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true})
}

func deleteScript(c echo.Context) error {
	script := models.Script{}
	if res := models.DB.Delete(&script, c.Param("id")); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true})
}
