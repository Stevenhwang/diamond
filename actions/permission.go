package actions

// import (
// 	"diamond/models"
// 	"net/http"

// 	"github.com/labstack/echo/v4"
// )

// func getPermissions(c echo.Context) error {
// 	var total int64
// 	permissions := models.Permissions{}
// 	models.DB.Model(&models.Permission{}).Scopes(models.Filter(models.Permission{}, c)).Count(&total)
// 	result := models.DB.Scopes(models.Filter(models.Permission{}, c), models.Paginate(c)).Find(&permissions)
// 	if result.Error != nil {
// 		return c.JSON(http.StatusOK, H{"code": 1, "message": result.Error.Error()})
// 	}
// 	return c.JSON(http.StatusOK, H{"code": 0, "data": permissions, "total": total})
// }

// func createPermission(c echo.Context) error {
// 	permission := &models.Permission{}
// 	if err := c.Bind(permission); err != nil {
// 		return c.JSON(http.StatusOK, H{"code": 1, "message": err.Error()})
// 	}
// 	if err := c.Validate(permission); err != nil {
// 		return c.JSON(http.StatusOK, H{"code": 2, "message": err.Error()})
// 	}
// 	if result := models.DB.Create(permission); result.Error != nil {
// 		return c.JSON(http.StatusOK, H{"code": 3, "message": result.Error.Error()})
// 	}
// 	return c.JSON(http.StatusOK, H{"code": 0, "message": "create permission success"})
// }

// func updatePermission(c echo.Context) error {
// 	permission := &models.Permission{}
// 	if result := models.DB.Find(permission, c.Param("id")); result.Error != nil {
// 		return c.JSON(http.StatusOK, H{"code": 1, "message": result.Error.Error()})
// 	}
// 	if err := c.Bind(permission); err != nil {
// 		return c.JSON(http.StatusOK, H{"code": 2, "message": err.Error()})
// 	}
// 	if err := c.Validate(permission); err != nil {
// 		return c.JSON(http.StatusOK, H{"code": 3, "message": err.Error()})
// 	}
// 	if result := models.DB.Select("*").Updates(permission); result.Error != nil {
// 		return c.JSON(http.StatusOK, H{"code": 4, "message": result.Error.Error()})
// 	}
// 	return c.JSON(http.StatusOK, H{"code": 0, "message": "update permission success"})
// }

// func deletePermission(c echo.Context) error {
// 	permission := &models.Permission{}
// 	if result := models.DB.Delete(permission, c.Param("id")); result.Error != nil {
// 		return c.JSON(http.StatusOK, H{"code": 1, "message": result.Error.Error()})
// 	}
// 	return c.JSON(http.StatusOK, H{"code": 0, "message": "delete permission success"})
// }
