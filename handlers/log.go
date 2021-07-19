package handlers

import (
	"diamond/models"

	"github.com/gin-gonic/gin"
)

// 获取系统日志
func LogListPerm(c *gin.Context) {
	logs, total, err := models.GetLogList(c)
	if err != nil {
		respMsg(c, 1, err.Error())
	}
	respData(c, 0, logs, total)
}
