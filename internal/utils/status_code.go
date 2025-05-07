package utils

/**
 * @File: status_code.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/14 下午3:34
 * @Software: GoLand
 * @Version:  1.0
 */

// Success
// 表示成功狀態
const Success = "0000"

// ErrorCode
// 表示錯誤狀態
type ErrorCode struct {
	StatusCode string
	Message    string
}

var (
	CodeBadRequest   = ErrorCode{"4000", "Bad Request: Invalid format"}
	CodeParamInvalid = ErrorCode{"4001", "Invalid parameters"}
	CodeEmailExists  = ErrorCode{"4002", "Email already exists"}
	CodeUnauthorized = ErrorCode{"4010", "Unauthorized"}
	CodeServerError  = ErrorCode{"5000", "Internal server error"}
)
