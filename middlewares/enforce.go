package middlewares

import (
	"diamond/policy"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Enforce 中间件使用casbin进行权限检查
func Enforce(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Get("username").(string)
		pass, err := policy.Enforcer.Enforce(username, c.Request().URL.Path, c.Request().Method)
		if err != nil {
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		}
		if !pass {
			return echo.NewHTTPError(http.StatusForbidden, "This action is forbidden")
		}
		return next(c)
	}
}
