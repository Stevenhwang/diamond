package models

import (
	"diamond/misc"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func init() {
	host := misc.Config.GetString("mysql.host")
	port := misc.Config.GetInt("mysql.port")
	user := misc.Config.GetString("mysql.user")
	password := misc.Config.GetString("mysql.password")
	dbName := misc.Config.GetString("mysql.dbName")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", user, password, host, port, dbName)
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		misc.Logger.Fatal().Err(err).Str("from", "db").Msg("connect mysql error")
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(5)
	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(20)
	// SetConnMaxLifetime 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)
	// 迁移 schema
	db.AutoMigrate(&User{}, &Server{}, &Credential{}, &Permission{}, &Record{}, &Task{}, &TaskHistory{})
	// seed admin user
	// var count int64
	// db.Model(&User{}).Where("username = ?", "admin").Count(&count)
	// if count == 0 {
	// 	admin := User{Username: "admin", Password: "12345678", IsActive: true}
	// 	db.Create(&admin)
	// }
	DB = db
}
