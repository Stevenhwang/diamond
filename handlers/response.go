package handlers

import (
	"github.com/gin-gonic/gin"
)

func respMsg(c *gin.Context, code int, message string) {
	c.JSON(200, gin.H{
		"code":    code,
		"message": message,
	})
}

func respData(c *gin.Context, code int, data interface{}, total int64) {
	c.JSON(200, gin.H{
		"code":  code,
		"data":  data,
		"total": total,
	})
}
