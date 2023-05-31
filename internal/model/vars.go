package model

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DbEngin *gorm.DB

func init() {
	db, err := gorm.Open(sqlite.Open("codetpl.db"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// 设置数据库连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	db.AutoMigrate(&Ds{}, &Table{}, &Column{})
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	DbEngin = db.Debug()
}
