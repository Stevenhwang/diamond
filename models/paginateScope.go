package models

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 分页
func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(c.Query("page"))
		if page <= 0 {
			page = 1
		}
		pageSize, _ := strconv.Atoi(c.Query("limit"))
		if pageSize <= 0 {
			pageSize = 15
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
