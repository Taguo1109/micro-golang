package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"micro-golang/internal/config"
	"micro-golang/internal/middlewares"
	"micro-golang/internal/order"
	"os"
)

/**
 * @File: main.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/23 上午11:11
 * @Software: GoLand
 * @Version:  1.0
 */

func main() {

	// 新增log 完整資訊
	initializeLogger()

	// DB初始化
	config.ConnectDB()
	// Redis 初始化
	config.InitRedis()

	port := os.Getenv("ORDER_PORT")
	if port == "" {
		port = "9000"
	}

	// User services URL from env
	userSvcURL := os.Getenv("USER_SVC_URL")
	if userSvcURL == "" {
		userSvcURL = "http://localhost:8000"
	}

	r := gin.Default()
	r.Use(middlewares.JWTAuth())
	oh := order.NewHandler(userSvcURL)
	r.GET("/orders/:id", oh.GetOrder)
	r.GET("/orders/email/:id", oh.GetOrderWithEmail)

	log.Printf("Order services running on :%s\n", port)
	log.Fatal(r.Run(":" + port))
}

// initializeLogger 增加 log 完整資訊
func initializeLogger() {
	// 結合標準旗標 (日期時間) 與短檔案名/行號
	//log.SetFlags(log.LstdFlags | log.Lshortfile)
	// 如果你想使用完整路徑，可以替換為：
	log.SetFlags(log.LstdFlags | log.Llongfile)

	// 你也可以將日誌輸出到檔案而不是標準錯誤輸出
	// file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err == nil {
	//     log.SetOutput(file)
	// } else {
	//     log.Println("無法開啟日誌檔案:", err)
	// }
}
