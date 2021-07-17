package handlers

import (
	"database/sql"
	"diamond/models"
	"diamond/utils.go"

	"github.com/gofiber/fiber/v2"
	fiberUtils "github.com/gofiber/fiber/v2/utils"
)

func authTokenMW(c *fiber.Ctx) error {
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
			"code":    1,
			"message": "Token非法!",
		})
	}
	rdToken := utils.GetToken(uid)
	if len(rdToken) == 0 || rdToken != xToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    1,
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
