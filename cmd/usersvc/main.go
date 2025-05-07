package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"micro-golang/internal/config"
	"micro-golang/internal/middlewares"
	"micro-golang/internal/user"
	"net/http"
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

	port := os.Getenv("USER_PORT")
	if port == "" {
		port = "8000"
	}

	r := gin.Default()
	r.Use(middlewares.JWTAuth())
	uh := user.NewHandler()
	r.GET("/users/:id", uh.GetUser)
	r.GET("/users/email/:id", uh.GetUserEmail)
	r.GET("/api/v1/users/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"test": "測試nginx新的url",
		})
	})

	log.Printf("User services running on :%s\n", port)
	log.Fatal(r.Run(":" + port))

}
