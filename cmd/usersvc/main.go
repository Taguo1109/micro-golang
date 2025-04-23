package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"micro-golang/internal/middleware"
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

	port := os.Getenv("USER_PORT")
	if port == "" {
		port = "8000"
	}

	r := gin.Default()
	r.Use(middleware.JWTAuth())
	uh := user.NewHandler()
	r.GET("/users/:id", uh.GetUser)
	r.GET("/api/v1/users/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"test": "測試nginx新的url",
		})
	})

	log.Printf("User service running on :%s\n", port)
	log.Fatal(r.Run(":" + port))

}
