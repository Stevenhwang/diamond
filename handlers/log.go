package handlers

import (
	"diamond/models"

	"github.com/gofiber/fiber/v2"
)

// 获取系统日志
func LogListPerm(c *fiber.Ctx) error {
	logs, total, err := models.GetLogList(c)
	if err != nil {
		return RespMsgSuccess(c, 1, err.Error())
	}
	return RespDataSuccess(c, 0, logs, total)
}
