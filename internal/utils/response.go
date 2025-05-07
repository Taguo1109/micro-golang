package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * @File: response.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/14 下午2:55
 * @Software: GoLand
 * @Version:  1.0
 */

type JsonResult struct {
	StatusCode string      `json:"status_code"`
	Msg        interface{} `json:"msg"`
	MsgDetail  string      `json:"msg_detail"`
	Data       interface{} `json:"data"`
}

// ReturnSuccess 回傳成功
func ReturnSuccess(c *gin.Context, data interface{}, detailMsg ...string) {
	response := JsonResult{
		StatusCode: Success,
		Msg:        "Success",
		Data:       data,
	}
	if len(detailMsg) > 0 {
		response.MsgDetail = detailMsg[0]
	}

	c.JSON(http.StatusOK, response)
}

// ReturnError 回傳統一格式的錯誤 JSON 響應。
// 參數 	code 為業務錯誤碼，detailMsg 為可選的詳細錯誤說明（msg_detail）。
func ReturnError(c *gin.Context, errCode ErrorCode, data interface{}, detailMsg ...string) {
	response := JsonResult{
		StatusCode: errCode.StatusCode,
		Msg:        errCode.Message,
		Data:       data,
	}
	if len(detailMsg) > 0 {
		response.MsgDetail = detailMsg[0]
	}

	c.JSON(http.StatusOK, response)
}
