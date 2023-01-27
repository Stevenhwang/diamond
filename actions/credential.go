package actions

import (
	"diamond/models"

	"github.com/labstack/echo/v4"
)

func getCredentials(c echo.Context) error {
	var total int64
	credentials := models.Credentials{}
	if res := models.DB.Model(&models.Credential{}).Count(&total); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	if res := models.DB.Scopes(models.Paginate(c)).Omit("auth_content").Find(&credentials); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true, "data": credentials, "total": total})
}

func createCredential(c echo.Context) error {
	credential := models.Credential{}
	if err := c.Bind(&credential); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if err := c.Validate(&credential); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if result := models.DB.Create(&credential); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true})
}

func updateCredential(c echo.Context) error {
	credential := models.Credential{}
	if res := models.DB.First(&credential, c.Param("id")); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	if err := c.Bind(&credential); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	if err := c.Validate(&credential); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	// 处理auth_content更新
	excludeColumns := []string{}
	if len(credential.AuthContent) == 0 {
		excludeColumns = append(excludeColumns, "auth_content")
	}
	if result := models.DB.Select("*").Omit(excludeColumns...).Updates(&credential); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true})
}

func deleteCredential(c echo.Context) error {
	credential := models.Credential{}
	if result := models.DB.First(&credential, c.Param("id")); result.Error != nil {
		return echo.NewHTTPError(400, result.Error.Error())
	}
	// ensure delete hook run
	if res := models.DB.Delete(&credential); res.Error != nil {
		return echo.NewHTTPError(400, res.Error.Error())
	}
	return c.JSON(200, echo.Map{"success": true})
}
