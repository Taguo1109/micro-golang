package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"micro-golang/internal/config"
	"micro-golang/internal/middlewares"
	"micro-golang/internal/user"
	"net/http"
	"os"
	"time"
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
	
	// 需權限的名單
	r.Use(middlewares.JWTAuth())
	// 加入全域錯誤攔截器
	r.Use(middlewares.GlobalErrorHandler())
	// 跨域設定
	setupCorsMiddleware(r)

	uh := user.NewHandler()
	r.GET("/users/:id", uh.GetUser)
	r.GET("/users/email/:id", uh.GetUserEmail)
	r.GET("/users/profile", uh.GetProfile)
	r.GET("/api/v1/users/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"test": "測試nginx新的url",
		})
	})

	log.Printf("User services running on :%s\n", port)
	log.Fatal(r.Run(":" + port))

}

// setupCorsMiddleware 設置CORS中間件
func setupCorsMiddleware(r *gin.Engine) {
	// 加入 CORS 設定
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",       // 本地開發用
			"https://taguo1109.github.io", // GitHub Pages 正式站
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // 如果你有用 cookie/token
		MaxAge:           12 * time.Hour,
	}))
}
