package models

import (
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// 数据过滤，链式，like查询，针对string字段
func Filter(model interface{}, c *fiber.Ctx) func(db *gorm.DB) *gorm.DB {
	reflectType := reflect.ValueOf(model).Type()
	return func(db *gorm.DB) *gorm.DB {
		for i := 0; i < reflectType.NumField(); i++ {
			field := reflectType.Field(i).Tag.Get("filter")
			if f := c.Query(field); len(f) > 0 {
				qString := fmt.Sprintf("%s like ?", field)
				db.Where(qString, "%"+f+"%")
			}
		}
		return db
	}
}
