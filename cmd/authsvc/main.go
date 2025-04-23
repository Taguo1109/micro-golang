package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	r := gin.Default()
	r.POST("/login", loginHandler)
	r.Run(":7000")
}

type loginReq struct{ User, Pass string }
type loginResp struct{ Token string }

func loginHandler(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "格式錯誤"})
		return
	}
	// TODO: 在這裡驗 credentials (DB or in-memory)
	if !(req.User == "demo" && req.Pass == "1234") {
		c.JSON(401, gin.H{"error": "帳密錯誤"})
		return
	}

	// 簽發 JWT
	claims := jwt.RegisteredClaims{
		Subject:   req.User,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokStr, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(500, gin.H{"error": "無法簽發 Token"})
		return
	}
	c.JSON(200, loginResp{Token: tokStr})
}
