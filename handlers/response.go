package handlers

import "github.com/gofiber/fiber/v2"

func RespMsgSuccess(c *fiber.Ctx, code int, message string) error {
	return c.JSON(fiber.Map{
		"code":    code,
		"message": message,
	})
}

func RespDataSuccess(c *fiber.Ctx, code int, data interface{}, total int64) error {
	return c.JSON(fiber.Map{
		"code":  code,
		"data":  data,
		"total": total,
	})
}

func RespError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"code":    500,
		"message": message,
	})
}
