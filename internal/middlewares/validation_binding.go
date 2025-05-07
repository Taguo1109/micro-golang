package middlewares

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

/**
 * @File: validation_binding.go
 * @Description:
 *
 * 註冊驗證資訊
 *
 * @Author: Timmy
 * @Create: 2025/4/15 下午5:31
 * @Software: GoLand
 * @Version:  1.0
 */

// UserPwd 密碼驗證
// 必須要有一個大寫英文一個小寫英文及數字，至少6碼最多30碼
// Go 不支援 Lookahead 語法(?= 之類的 所以 ^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{6,30}$ 無法
func UserPwd(field validator.FieldLevel) bool {
	password := field.Field().String()

	if len(password) < 6 || len(password) > 30 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	isOnlyAlphaNum := regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(password)

	return hasUpper && hasLower && hasDigit && isOnlyAlphaNum
}

// UserName 用戶名稱驗證
// 只能英文數字，至少6字最多20字
func UserName(field validator.FieldLevel) bool {
	if match, _ := regexp.MatchString(`^[a-zA-Z0-9]{6,20}$`, field.Field().String()); match {
		return true
	}
	return false
}
