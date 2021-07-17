package models

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Log struct {
	ID        uint
	Username  string         `gorm:"size:128" json:"username"`
	IP        string         `gorm:"size:128" json:"ip"`
	Method    string         `gorm:"size:16" json:"method"`
	URL       string         `gorm:"size:128" json:"url"`
	Data      sql.NullString `gorm:"type:text" json:"data"`
	CreatedAt time.Time      `json:"created_at"`
}

type Logs []Log

func GetLogList(c *fiber.Ctx) (logs Logs, total int64, err error) {
	logs = Logs{}
	query := DB.Scopes(Filter(Log{}, c))
	if len(c.Params("date_before")) > 0 {
		query.Where("created_at <= ?", c.Params("date_before"))
	}
	if len(c.Params("date_after")) > 0 {
		query.Where("created_at >= ?", c.Params("date_after"))
	}
	query.Count(&total)
	result := query.Scopes(Paginate(c)).Order("created_at desc").Find(&logs)
	return logs, total, result.Error
}
