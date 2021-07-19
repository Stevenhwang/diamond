package handlers

import (
	"database/sql"
	"diamond/models"
	"diamond/utils.go"

	"github.com/gofiber/fiber/v2"
	fiberUtils "github.com/gofiber/fiber/v2/utils"
)

func realIPMW(c *fiber.Ctx) error {
	if xRealIP := c.Get("X-Real-Ip"); len(xRealIP) > 0 {
		c.Locals("user_ip", xRealIP)
		return c.Next()
	}
	if xForwardedFor := c.IPs(); len(xForwardedFor) > 0 {
		c.Locals("user_ip", xForwardedFor[0])
		return c.Next()
	}
	c.Locals("user_ip", c.IP())
	return c.Next()
}

func authTokenMW(c *fiber.Ctx) error {
	// 登录接口跳过
	if c.Path() == "/login" {
		return c.Next()
	}
	// 检查token
	xToken := c.Get("X-Token")
	if len(xToken) == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    1,
			"message": "请登录!",
		})
	}
	uid, username, isSuperuser := utils.J.DecodeToken(xToken)
	if uid == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    2,
			"message": "Token非法!",
		})
	}
	rdToken := utils.GetToken(uid)
	if len(rdToken) == 0 || rdToken != xToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    3,
			"message": "Token失效!",
		})
	}
	c.Locals("uid", uid)
	c.Locals("username", username)
	c.Locals("is_superuser", isSuperuser)
	// 记录访问日志(除get请求)
	if c.Method() != "GET" {
		log := &models.Log{}
		log.Username = username
		log.IP = c.IPs()[0]
		log.Method = c.Method()
		log.URL = c.Path()
		b := fiberUtils.CopyString(string(c.Body()))
		if len(b) > 0 {
			log.Data = sql.NullString{String: b, Valid: true}
		}
		go func(log *models.Log) {
			models.DB.Create(log)
		}(log)
	}
	return c.Next()
}
