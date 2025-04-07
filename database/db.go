package database

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// DB 全局数据库连接实例
var DB *gorm.DB

// Initialize 初始化数据库连接
func Initialize(dsn string) error {
	var err error
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// 开启连接池
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)

	// 启用日志
	DB.LogMode(true)

	log.Println("数据库连接成功")
	return nil
}

// Close 关闭数据库连接
func Close() {
	if DB != nil {
		DB.Close()
	}
}
