package middleware

import (
	"github.com/gin-gonic/gin"
	"micro-golang/pkg/jwt"
	"net/http"
	"strings"
)

/**
 * @File: auth.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/23 下午3:02
 * @Software: GoLand
 * @Version:  1.0
 */

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "請帶上 Bearer Token"})
			return
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")
		claims, err := jwt.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token 無效或過期"})
			return
		}
		// 可把 user id or subject 存進 Context
		c.Set("user", claims.Subject)
		c.Next()
	}
}
