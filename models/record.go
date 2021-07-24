package models

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Record struct {
	ID        uint      `json:"id"`
	User      string    `gorm:"size:128" json:"user" filter:"user"`
	ServerID  uint      `json:"server_id"`
	CreatedAt time.Time `json:"created_at"`
	File      string    `gorm:"size:128" json:"file"`
}

type Records []Record

func (r *Record) AfterDelete(tx *gorm.DB) (err error) {
	if err := os.Remove(r.File); err != nil {
		log.Println("already removed")
	}
	log.Println("removed")
	return
}

func GetRecordList(c *gin.Context) (records Records, total int64, err error) {
	records = Records{}
	query := DB.Model(&Record{}).Scopes(Filter(Record{}, c))
	if len(c.Query("date_before")) > 0 {
		query.Where("created_at <= ?", c.Query("date_before"))
	}
	if len(c.Query("date_after")) > 0 {
		query.Where("created_at >= ?", c.Query("date_after"))
	}
	query.Count(&total)
	result := query.Scopes(Paginate(c)).Order("created_at desc").Find(&records)
	return records, total, result.Error
}
