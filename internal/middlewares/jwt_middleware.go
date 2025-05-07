package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"micro-golang/internal/config"
	"micro-golang/internal/utils"
	"net/http"
	"strings"
)

/**
 * @File: jwt_middleware.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/23 下午3:02
 * @Software: GoLand
 * @Version:  1.0
 */

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 從 Header 中讀取 Authorization 欄位
		authHeader := c.GetHeader("Authorization")

		// 2. 檢查 Header 是否以 "Bearer " 開頭
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, utils.JsonResult{
				StatusCode: "401",
				Msg:        "No token provided in Authorization header",
				MsgDetail:  "請先登入或確認 Request Header 中的 Authorization 格式是否為 Bearer token",
			})
			c.Abort()
			return
		}

		// 3. 擷取 Bearer Token 的實際內容
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 4. Redis 黑名單檢查
		blacklisted, _ := config.RDB.Exists(config.Ctx, "blacklist:access_token:"+tokenString).Result()
		if blacklisted == 1 {
			c.JSON(http.StatusUnauthorized, utils.JsonResult{
				StatusCode: "401",
				Msg:        "Token is Logout and inValid",
				MsgDetail:  "Token 已被登出或無效",
			})
			c.Abort()
			return
		}

		// 5. 驗證並解析 JWT Token
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return utils.JwtKey, nil
		})

		// 6. 確認 token 有效，若驗證失敗則中止請求
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, utils.JsonResult{
				StatusCode: "401",
				Msg:        "Invalid token",
				MsgDetail:  "Token 無效或已過期，請重新登入",
			})
			c.Abort()
			return
		}

		// 🚨 token_type 檢查：若為 refresh，拒絕使用
		if tokenType, ok := claims["token_type"]; ok && tokenType == "refresh" {
			c.JSON(http.StatusUnauthorized, utils.JsonResult{
				StatusCode: "401",
				Msg:        "Invalid token type",
				MsgDetail:  "請使用 access token 進行此操作",
			})
			c.Abort()
			return
		}

		// 7. 從 claims 中取出使用者資訊，設定到 Context 讓後續 handlers 使用
		c.Set("email", claims["email"])
		c.Set("userId", claims["userId"])
		c.Set("role", claims["role"])

		// 8. 放行
		c.Next()
	}
}
