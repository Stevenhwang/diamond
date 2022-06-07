package models

import (
	"log"
	"os"
	"time"

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
