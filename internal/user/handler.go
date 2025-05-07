package user

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"micro-golang/internal/config"
	"micro-golang/internal/dto"
	"micro-golang/internal/models"
	"micro-golang/internal/utils"
	"net/http"
	"time"
)

/**
 * @File: handler.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/23 上午11:11
 * @Software: GoLand
 * @Version:  1.0
 */

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

func (h *Handler) GetUser(c *gin.Context) {
	// 模擬 DB
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id":   id,
		"name": "User" + id,
	})
}

func (h *Handler) GetUserEmail(c *gin.Context) {
	// 模擬 DB
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"email": id + "@example.com",
	})
}

// GetProfile 獲取用戶基本資料
func (h *Handler) GetProfile(c *gin.Context) {
	emailVal, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No email in token"})
		return
	}
	email := emailVal.(string)

	// 1️⃣ 先從 Redis 查快取
	cacheKey := "user:" + email
	cached, err := config.RDB.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		// 如果有快取，直接回傳
		var cachedUser dto.UserLoginResponseDTO
		// json.Unmarshal 將資料JSON格式化
		if err := json.Unmarshal([]byte(cached), &cachedUser); err == nil {
			utils.ReturnSuccess(c, cachedUser, "from cache")
			return
		}
	}

	// 2️⃣ 沒快取，查 DB
	var user models.User
	result := config.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 3️⃣ 查到後，存入 Redis 快取（設 10 分鐘過期）
	safeUser := dto.UserLoginResponseDTO{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Role:     user.Role,
	}
	userBytes, _ := json.Marshal(safeUser)
	config.RDB.Set(config.Ctx, cacheKey, userBytes, 10*time.Minute)
	utils.ReturnSuccess(c, safeUser, "from db")
}
