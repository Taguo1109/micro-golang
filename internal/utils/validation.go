package utils

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

/**
 * @File: validation.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/15 下午5:06
 * @Software: GoLand
 * @Version:  1.0
 */

// ExtractFieldErrorMessages 根據傳入 struct 的欄位 tag `validateMsg`
// 對 validator 驗證錯誤逐一解析，並組成一份 map[string]string 的欄位錯誤訊息集合。
// 欄位名稱會轉換為對應的 `json` tag（如 "Username" -> "username"）。
func ExtractFieldErrorMessages(obj interface{}, errs validator.ValidationErrors) map[string]string {
	errMap := make(map[string]string)
	objT := reflect.TypeOf(obj)

	// 如果是指標類型（pointer），取得實際結構體 type
	if objT.Kind() == reflect.Ptr {
		objT = objT.Elem()
	}

	// 逐筆驗證錯誤處理
	for _, fieldErr := range errs {
		fieldName := fieldErr.Field() // Go 的欄位名稱（如 Username）
		jsonTag := fieldName          // 預設用 Go 欄位名

		// 透過 reflect 找到 struct 欄位定義
		if f, found := objT.FieldByName(fieldName); found {
			// 嘗試取得 json 標籤欄位名
			if tag := f.Tag.Get("json"); tag != "" {
				jsonTag = strings.Split(tag, ",")[0] // 避免包含 ,omitempty
			}

			// 取得 validateMsg 自訂錯誤訊息的 tag
			tagMsg := f.Tag.Get("validateMsg")
			msgMap := parseValidateMsgTag(tagMsg)

			// 如果該驗證規則有自訂訊息，則使用
			if msg, ok := msgMap[fieldErr.Tag()]; ok {
				errMap[jsonTag] = msg
			} else {
				// 否則給預設訊息
				errMap[jsonTag] = "欄位格式錯誤：" + fieldErr.Tag()
			}
		}
	}

	return errMap
}

// parseValidateMsgTag 將 tag 字串解析為 map，例如：
// "required=必填,email=格式錯誤" → map["required"] = "必填", map["email"] = "格式錯誤"
func parseValidateMsgTag(tag string) map[string]string {
	result := make(map[string]string)
	pairs := strings.Split(tag, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			result[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return result
}
