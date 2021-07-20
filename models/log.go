package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/nulls"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Log struct {
	ID        uuid.UUID    `gorm:"type:char(36);primary_key" json:"id"`
	Username  string       `gorm:"size:128" json:"username"`
	IP        string       `gorm:"size:128" json:"ip"`
	Method    string       `gorm:"size:16" json:"method"`
	URL       string       `gorm:"size:128" json:"url"`
	Data      nulls.String `gorm:"type:text" json:"data"`
	CreatedAt time.Time    `json:"created_at"`
}

type Logs []Log

func (l *Log) BeforeCreate(tx *gorm.DB) (err error) {
	x, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	l.ID = x
	return nil
}

func GetLogList(c *gin.Context) (logs Logs, total int64, err error) {
	logs = Logs{}
	query := DB.Model(&Log{}).Scopes(Filter(Log{}, c))
	if len(c.Query("date_before")) > 0 {
		query.Where("created_at <= ?", c.Query("date_before"))
	}
	if len(c.Query("date_after")) > 0 {
		query.Where("created_at >= ?", c.Query("date_after"))
	}
	query.Count(&total)
	result := query.Scopes(Paginate(c)).Order("created_at desc").Find(&logs)
	return logs, total, result.Error
}
