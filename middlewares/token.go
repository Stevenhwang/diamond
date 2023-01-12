package middlewares

import (
	"diamond/models"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// Token 中间件用来解析 jwt token
func Token(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return echo.NewHTTPError(400, "JWT token missing or invalid")
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return echo.NewHTTPError(400, "failed to cast claims as jwt.MapClaims")
		}
		uid := uint(claims["uid"].(float64))
		username := claims["username"].(string)
		// 查看用户是否被禁用，有可能是登录之后被禁用，所以需要查询实时数据
		user := models.User{}
		if res := models.DB.First(&user, uid); res.Error != nil {
			return echo.NewHTTPError(400, res.Error.Error())
		}
		if !user.IsActive {
			return echo.NewHTTPError(401, "账号禁用")
		}
		// set context
		c.Set("uid", uid)
		c.Set("username", username)
		c.Set("menus", user.Menus)
		return next(c)
	}
}
