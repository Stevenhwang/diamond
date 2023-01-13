package middlewares

import (
	"diamond/cache"
	"net/http"

	"github.com/labstack/echo/v4"
)

// BanIP 中间件检查IP是否在黑名单
func BanIP(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Echo/Golang")
		b, err := cache.GetBan(c.RealIP())
		if err != nil {
			return echo.NewHTTPError(400, err.Error())
		}
		if b {
			return echo.NewHTTPError(http.StatusForbidden, "ip forbiden")
		}
		return next(c)
	}
}
