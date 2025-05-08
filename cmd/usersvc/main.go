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

	// 跨域設定(要先放跨域)在擺權限過濾
	setupCorsMiddleware(r)
	// 需權限的名單
	r.Use(middlewares.JWTAuth())
	// 加入全域錯誤攔截器
	r.Use(middlewares.GlobalErrorHandler())

	// 1. 獲取 DB 和 Redis 客戶端實例
	// 這些應該在你的應用程式生命週期的早期就被初始化。
	// 例如，如果它們可以通過 config 套件全局訪問：
	dbInstance := config.DB           // 你的 *gorm.DB 實例
	redisClientInstance := config.RDB // 你的 *redis.Client 實例
	// 2. 創建 user.Service 的實例
	// NewService 是在你的 user 套件中定義的
	userServiceInstance := user.NewService(dbInstance, redisClientInstance)
	// 3. 創建 user.Handler 的實例，並傳入 userServiceInstance
	// NewHandler 也是在你的 user 套件中定義的
	uh := user.NewHandler(userServiceInstance)

	ur := r.Group("/users")
	ur.GET("/:id", uh.GetUser)
	ur.GET("/email/:id", uh.GetUserEmail)
	// 獲取個人資料
	ur.GET("/profile", uh.GetProfile)
	// 更新個人資料
	ur.PUT("/profile", uh.UpdateProfile)

	ur.GET("/api/v1/users/:id", func(c *gin.Context) {
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
