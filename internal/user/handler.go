package user

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"micro-golang/internal/config"
	"micro-golang/internal/dto"
	"micro-golang/internal/models"
	"micro-golang/internal/utils"
	"net/http"
	"strings"
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

type Handler struct {
	userService *Service // 添加 userService 字段
}

// NewHandler 修改構造函數以接收 UserService
func NewHandler(userService *Service) *Handler {
	return &Handler{userService: userService}
}

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

// UpdateProfile 更新用戶基本資料 (Username, Email) - 調用 Service
func (h *Handler) UpdateProfile(c *gin.Context) {
	emailVal, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No email in token"})
		return
	}
	currentEmail := emailVal.(string)

	var req dto.UserUpdateProfileDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// 調用 service 方法
	// 使用 c.Request.Context() 將請求上下文傳遞給 service 層，這對於超時控制和值傳遞很有用
	updatedUserDTO, err := h.userService.UpdateUserProfile(c.Request.Context(), currentEmail, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, ErrEmailInUse):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, ErrUpdateNoChanges):
			// 如果 ErrUpdateNoChanges 伴隨 updatedUserDTO 返回，說明請求有效但無實際變更
			log.Println(err.Error())
			log.Println(err)
			if updatedUserDTO != nil {
				utils.ReturnSuccess(c, updatedUserDTO, "No effective changes to user profile.")
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "No fields provided for update or no changes detected"})
			}
		case strings.Contains(err.Error(), "validation failed"): // 簡易判斷，可以定義更精確的 error
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case strings.Contains(err.Error(), "cache marshal failed") || strings.Contains(err.Error(), "cache set failed"):
			// 操作成功，但快取有問題，仍然返回成功，但可以記錄日誌
			// log.Warn("UpdateProfile successful but cache operation failed: %v", err)
			utils.ReturnSuccess(c, updatedUserDTO, err.Error()) // updatedUserDTO 應該會有值
		default:
			// 其他未知錯誤或數據庫錯誤
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile: " + err.Error()})
		}
		return
	}

	utils.ReturnSuccess(c, updatedUserDTO, "User profile updated successfully")
}
