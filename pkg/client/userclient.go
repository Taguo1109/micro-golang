package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/**
 * @File: userclient.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/23 上午11:12
 * @Software: GoLand
 * @Version:  1.0
 */

type UserClient struct {
	baseURL string
}

func NewUserClient(baseURL string) *UserClient {
	return &UserClient{baseURL: baseURL}
}

// 通用的 fetchData 函式，接收 http 方法作為參數
func (uc *UserClient) fetchData(method string, urlPath string, token string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s%s", uc.baseURL, urlPath)
	req, err := http.NewRequest(method, url, nil) // 使用傳入的 method
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf(token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("關閉 response Body 時發生錯誤:", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("請求失敗，狀態碼: %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func (uc *UserClient) FetchUser(id string, token string) (map[string]interface{}, error) {
	return uc.fetchData("GET", fmt.Sprintf("/users/%s", id), token)
}

func (uc *UserClient) FetchUserEmail(id string, token string) (map[string]interface{}, error) {
	return uc.fetchData("GET", fmt.Sprintf("/users/email/%s", id), token)
}
