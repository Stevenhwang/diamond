package middlewares

import (
	"diamond/misc"
	"net/http"

	"github.com/labstack/echo/v4"
)

// BanIP 中间件检查IP是否在黑名单
func BanIP(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Echo/Golang")
		_, err := misc.Cache.Get(c.RealIP())
		if err == nil { // 找到黑名单IP
			return echo.NewHTTPError(http.StatusForbidden, "ip forbiden")
		}
		return next(c)
	}
}
