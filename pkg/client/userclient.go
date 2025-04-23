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

func (uc *UserClient) FetchUser(id string) (map[string]interface{}, error) {
	resp, err := http.Get(fmt.Sprintf("%s/users/%s", uc.baseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}
