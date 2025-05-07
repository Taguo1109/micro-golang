package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"micro-golang/internal/auth"
	"micro-golang/internal/config"
	"micro-golang/internal/middlewares"
	"time"
)

/**
 * @File: main.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/23 下午2:57
 * @Software: GoLand
 * @Version:  1.0
 */

var jwtKey = []byte("A9v$34Fjfl1.pMv@z6k1!qwe93Km!")

func main() {

	// DB初始化
	config.ConnectDB()
	// Redis 初始化
	config.InitRedis()

	// 驗證器設定
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 註冊密碼驗證器
		_ = v.RegisterValidation("pwd_validation", middlewares.UserPwd)
		// 註冊使用者名稱驗證器
		_ = v.RegisterValidation("username_validation", middlewares.UserName)
	}

	r := gin.Default()
	// 加入全域錯誤攔截器
	r.Use(middlewares.GlobalErrorHandler())

	// 跨域設定
	setupCorsMiddleware(r)

	authGroup := r.Group("/auth")
	ah := auth.NewHandler(auth.Service{})
	authGroup.POST("/login", ah.Login)
	authGroup.POST("/register", ah.Register)
	authGroup.POST("/refresh", ah.RefreshToken)
	authGroup.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "測試是否自動部署"})
	})
	authGroup.POST("/logout", ah.LogoutHandler)
	err := r.Run(":7001")
	if err != nil {
		return
	}
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
