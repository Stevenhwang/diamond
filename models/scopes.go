package models

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// 分页
func Paginate(c echo.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page <= 0 {
			page = 1
		}
		pageSize, _ := strconv.Atoi(c.QueryParam("limit"))
		if pageSize <= 0 {
			pageSize = 15
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// and，like查询，带filter tag的string字段，需要查询端带相应字段来查询
func Filter(model interface{}, c echo.Context) func(db *gorm.DB) *gorm.DB {
	reflectType := reflect.ValueOf(model).Type()
	return func(db *gorm.DB) *gorm.DB {
		for i := 0; i < reflectType.NumField(); i++ {
			field := reflectType.Field(i).Tag.Get("filter")
			if f := c.QueryParam(field); len(f) > 0 {
				qString := fmt.Sprintf("%s like ?", field)
				db.Where(qString, "%"+f+"%")
			}
		}
		return db
	}
}

// or，like查询，带filter tag的string字段，不需要携带字段，只要带查询关键字query
func AnyFilter(model interface{}, c echo.Context) func(db *gorm.DB) *gorm.DB {
	query := c.QueryParam("query")
	if len(query) > 0 {
		reflectType := reflect.ValueOf(model).Type()
		return func(db *gorm.DB) *gorm.DB {
			for i := 0; i < reflectType.NumField(); i++ {
				if field := reflectType.Field(i).Tag.Get("filter"); len(field) > 0 {
					qString := fmt.Sprintf("%s like ?", field)
					db.Or(qString, "%"+query+"%")
				}
			}
			return db
		}
	} else {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}
}
