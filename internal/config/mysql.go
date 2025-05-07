package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

/**
 * @File: config.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/8 下午3:00
 * @Software: GoLand
 * @Version:  1.0
 */

var DB *gorm.DB

// ConnectDB 初始化 MySQL 資料庫連線，並加入重試機制避免 Docker 啟動順序導致錯誤。
// 在本地端通常不會遇到問題，但在 Docker Compose 中，go-app 可能在 MySQL 尚未準備好前就啟動，導致 connection refused。
// Retry up to 10 times, waiting 2 seconds each time.
//
// 💡 注意：若在 Docker compose 環境中未加 retry，可能會出現：
// dial tcp 172.xx.xx.xx:3306: connect: connection refused
func ConnectDB() {
	env := os.Getenv("APP_ENV") // 要去docker-compose.yml environment 加 APP_ENV
	var envFile string

	if env == "docker" {
		envFile = ".env"
	} else {
		envFile = ".env.local"
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Printf("⚠️ 無法載入 %s，將改用系統環境變數（可能導致 DB 或 Redis 無法連線）", envFile)
	}

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	var db *gorm.DB
	var err error
	for i := 1; i <= 10; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			DB = db
			log.Println("✅ Database connected!")
			return
		}
		log.Printf("❌ DB 連線失敗 (%d/10)：%v\n", i, err)
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("❌ 最終連線資料庫失敗：%v", err)
}
