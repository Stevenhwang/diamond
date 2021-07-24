package models

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 数据过滤，链式，like查询，针对string字段
func Filter(model interface{}, c *gin.Context) func(db *gorm.DB) *gorm.DB {
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

// 任意带filter标签的字段都过滤，or条件，like查询，针对string字段
func AnyFilter(model interface{}, query string) func(db *gorm.DB) *gorm.DB {
	reflectType := reflect.ValueOf(model).Type()
	return func(db *gorm.DB) *gorm.DB {
		for i := 0; i < reflectType.NumField(); i++ {
			field := reflectType.Field(i).Tag.Get("filter")
			qString := fmt.Sprintf("%s like ?", field)
			db.Or(qString, "%"+query+"%")
		}
		return db
	}
}
