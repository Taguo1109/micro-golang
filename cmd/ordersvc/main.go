package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"micro-golang/internal/middleware"
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
	port := os.Getenv("ORDER_PORT")
	if port == "" {
		port = "9000"
	}

	// User service URL from env
	userSvcURL := os.Getenv("USER_SVC_URL")
	if userSvcURL == "" {
		userSvcURL = "http://localhost:8000"
	}

	r := gin.Default()
	r.Use(middleware.JWTAuth())
	oh := order.NewHandler(userSvcURL)
	r.GET("/orders/:id", oh.GetOrder)

	log.Printf("Order service running on :%s\n", port)
	log.Fatal(r.Run(":" + port))
}
