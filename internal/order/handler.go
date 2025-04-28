package order

import (
	"github.com/gin-gonic/gin"
	"micro-golang/pkg/client"
	"net/http"
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
	uc *client.UserClient
}

func NewHandler(userSvcURL string) *Handler {
	return &Handler{
		uc: client.NewUserClient(userSvcURL),
	}
}

func (h *Handler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	// 模擬訂單資料
	order := map[string]interface{}{
		"id":     orderID,
		"item":   "Gadget",
		"amount": 99.9,
	}

	// 從 Order Service 的請求 Header 中獲取 Authorization Token
	token := c.GetHeader("Authorization")

	// 互打 User Service，傳遞 Token
	user, err := h.uc.FetchUser("123", token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order": order,
		"user":  user,
	})
}

func (h *Handler) GetOrderWithEmail(c *gin.Context) {
	orderID := c.Param("id")
	token := c.GetHeader("Authorization") // 假設從 Order Service 的請求 Header 中獲取 Token

	order := map[string]interface{}{
		"id":     orderID,
		"item":   "Another Product",
		"amount": 199.9,
	}

	userEmail, err := h.uc.FetchUserEmail("456", token) // 呼叫新的方法
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user email: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order":     order,
		"userEmail": userEmail,
	})
}
