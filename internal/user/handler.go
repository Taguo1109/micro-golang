package user

import (
	"github.com/gin-gonic/gin"
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
