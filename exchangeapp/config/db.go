package config

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 数据库连接配置
func InitDB() {

	dsn := AppConfig.Database.Dsn
	// 初始化数据库连接(db)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to initialize database,got error %v", err)
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(AppConfig.Database.MaxIdleConns) //设置数据库连接池的最大空闲连接数
	sqlDB.SetMaxOpenConns(AppConfig.Database.MaxOpenConns) //设置数据库连接池的最大连接数
	sqlDB.SetConnMaxLifetime(time.Hour)                    //设置数据库连接的最大生命周期为1小时

	if err != nil {
		log.Fatalf("Failed to configure database,got error %v", err)
	}

	db.AutoMigrate(
		&models.Article{},
		&models.Comment{},
		&models.User{},
		&models.InAppNotification{},
	)

	global.Db = db //将数据库连接赋值给全局变量Db
}
