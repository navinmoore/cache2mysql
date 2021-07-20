package mysql

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var dbLock sync.Mutex

func NewDB() *gorm.DB {
	dsn := "root:123456@tcp(127.0.0.1:3306)/cache2mysql?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("open error=", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("open error=", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		fmt.Println("ping error=", err)
	}

	// 设置空闲连接池中链接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// 设置打开数据库链接的最大数量
	sqlDB.SetMaxOpenConns(100)
	// 设置了链接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}

func MysqlInstance() *gorm.DB {
	if db != nil {
		return db
	}
	dbLock.Lock()
	defer dbLock.Unlock()
	if db != nil {
		return db
	}
	return NewDB()
}
