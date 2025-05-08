package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
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

	// 新增log 完整資訊
	initializeLogger()
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
