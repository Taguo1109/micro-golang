package auth

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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
 * @Create: 2025/5/7 下午2:50
 * @Software: GoLand
 * @Version:  1.0
 */

// Handler 負責處理認證相關的 HTTP 請求
type Handler struct {
	authService Service // 注入 AuthService
}

// NewHandler 是一個建構函式，用於建立 AuthHandler 的實例
// AuthService 應該在 main.go 或設定路由的地方被初始化並傳入
func NewHandler(authService Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

// Register 會員註冊
func (h *Handler) Register(c *gin.Context) {
	var input dto.UserRegisterDTO

	if err := c.ShouldBindJSON(&input); err != nil {

		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errFields := utils.ExtractFieldErrorMessages(input, ve)
			utils.ReturnError(c, utils.CodeParamInvalid, errFields, "欄位驗證失敗")
			return
		}
		utils.ReturnError(c, utils.CodeParamInvalid, err.Error())
		return
	}
	// 呼叫service，內部會直接回 JSON
	h.authService.Register(c, input)
}

// Login 會員登入
func (h *Handler) Login(c *gin.Context) {
	var input dto.UserLoginDTO
	var dbUser models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// 找出該 Email 的使用者
	result := config.DB.Where("email = ?", input.Email).First(&dbUser)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email does not exist"})
		return
	}

	// 檢查密碼
	err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
		return
	}

	// 產生 JWT token
	accessToken, refreshToken, err := utils.GenerateJWT(dbUser.Email, dbUser.ID, "User")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
	}

	// 快取使用者資料
	cacheKey := "user:" + dbUser.Email

	// 使dbUser 變成 JSON 標準格式 ，safeUser 存取需要的資訊進去
	safeUser := dto.UserDTO{
		ID:           dbUser.ID,
		Email:        dbUser.Email,
		Username:     dbUser.Username,
		Role:         dbUser.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	var responseDTO dto.UserLoginResponseDTO
	bytes, _ := json.Marshal(safeUser)
	err = json.Unmarshal(bytes, &responseDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.JsonResult{
			StatusCode: "500",
			Msg:        "Failed to unmarshal JSON",
			MsgDetail:  "JSON轉換失敗",
		})
		return
	}
	userBytes, _ := json.Marshal(responseDTO)
	config.RDB.Set(config.Ctx, cacheKey, userBytes, 10*time.Minute)

	utils.ReturnSuccess(c, safeUser, "Login successful")
}

// RefreshToken 重新獲取 Token
func (h *Handler) RefreshToken(c *gin.Context) {
	// 從 JSON 或 localStorage 帶進來
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.RefreshToken == "" {
		utils.ReturnError(c, utils.CodeParamInvalid, nil, "請提供 refresh_token")
		return
	}

	// 解析 token
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(input.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.JwtKey, nil
	})

	// 檢查Redis 是否存在黑名單
	isBlacklisted, _ := config.RDB.Exists(config.Ctx, "blacklist:refresh_token:"+input.RefreshToken).Result()
	if isBlacklisted == 1 {
		utils.ReturnError(c, utils.CodeUnauthorized, nil, "refresh_token 已失效，請重新登入")
		return
	}

	// 驗證 token 是否有效 & 是 refresh token
	if err != nil || !token.Valid || claims["token_type"] != "refresh" {
		utils.ReturnError(c, utils.CodeUnauthorized, nil, "refresh_token 無效或過期")
		return
	}

	email, ok := claims["email"].(string)
	if !ok {
		utils.ReturnError(c, utils.CodeUnauthorized, nil, "Token 內容無效")
		return
	}

	// 查詢使用者資料
	var dbUser models.User
	result := config.DB.Where("email = ?", email).First(&dbUser)
	if result.Error != nil {
		utils.ReturnError(c, utils.CodeUnauthorized, nil, "找不到使用者")
		return
	}

	// 產生新 token
	newAccessToken, newRefreshToken, err := utils.GenerateJWT(dbUser.Email, dbUser.ID, "User")
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.JsonResult{
			StatusCode: "500",
			Msg:        "Can't generate access token",
			MsgDetail:  "無法產生新 token",
		})
		return
	}

	// 統一回傳 DTO
	safeUser := dto.UserDTO{
		ID:           dbUser.ID,
		Email:        dbUser.Email,
		Username:     dbUser.Username,
		Role:         dbUser.Role,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}
	utils.ReturnSuccess(c, safeUser, "Token refreshed successfully")
}

// LogoutHandler 登出
func (h *Handler) LogoutHandler(c *gin.Context) {

	var input dto.UserDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ReturnError(c, utils.CodeBadRequest, nil, "格式錯誤")
		return
	}
	// 1️⃣ access_token 放進黑名單
	accessClaims, err := utils.ParseToken(input.AccessToken)
	if err == nil {
		if exp, ok := accessClaims["exp"].(float64); ok {
			ttl := time.Until(time.Unix(int64(exp), 0))
			// 若還沒過期則加入黑名單
			if ttl > 0 {
				config.RDB.Set(c, "blacklist:access_token:"+input.AccessToken, "1", ttl)
			}
		}
	}
	// 2️⃣ refresh_token 放進黑名單
	refreshClaims, err := utils.ParseToken(input.RefreshToken)
	if err == nil {
		if exp, ok := refreshClaims["exp"].(float64); ok {
			ttl := time.Until(time.Unix(int64(exp), 0))
			// 若還沒過期則加入黑名單
			if ttl > 0 {
				config.RDB.Set(c, "blacklist:refresh_token:"+input.RefreshToken, "1", ttl)
			}
		}
	}
	utils.ReturnSuccess(c, nil, "Logout successful")
}
