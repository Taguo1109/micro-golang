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
	// 互打 User Service
	user, err := h.uc.FetchUser("123")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order": order,
		"user":  user,
	})
}
