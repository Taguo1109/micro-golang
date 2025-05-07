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
