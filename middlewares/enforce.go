package middlewares

import (
	"diamond/policy"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Enforce 中间件使用casbin进行权限检查
func Enforce(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Echo/4.7")
		is_superuser := c.Get("is_superuser").(bool)
		uid := c.Get("uid").(int)
		if is_superuser {
			return next(c)
		} else {
			// 加前缀区分user和role
			sub := fmt.Sprintf("user::%d", uid)
			obj := fmt.Sprintf("%s %s", c.Request().Method, c.Request().URL.Path)
			pass, err := policy.Enforcer.Enforce(sub, obj, "route")
			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden, "This action is forbidden")
			}
			if !pass {
				return echo.NewHTTPError(http.StatusForbidden, "This action is forbidden")
			}
			return next(c)
		}
	}
}
