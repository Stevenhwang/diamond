package middlewares

import (
	"diamond/config"
	"diamond/utils"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// Token 中间件用来解析 jwt token 并检查是否合法
func Token(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user").(*jwt.Token)
		// get token from request
		t, _ := token.SignedString([]byte(config.Config.Get("jwt.secret").(string)))
		claims := token.Claims.(jwt.MapClaims)
		uid := uint(claims["uid"].(float64))
		username := claims["username"].(string)
		is_superuser := claims["is_superuser"].(bool)
		// set context
		c.Set("uid", uid)
		c.Set("username", username)
		c.Set("is_superuser", is_superuser)
		// get token from redis and compare
		rdToken := utils.GetToken(uid)
		if len(rdToken) == 0 || rdToken != t {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"code": 401, "message": "token invalid"})
		} else {
			return next(c)
		}
	}
}
