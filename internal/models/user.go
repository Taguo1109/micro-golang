package models

/**
 * @File: user.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/8 下午2:55
 * @Software: GoLand
 * @Version:  1.0
 */

import (
	"time"
)

// User 建立User Table
type User struct {
	ID        uint      `gorm:"primary" json:"id"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `json:"password"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 對應表名，若不加預設對應users
func (User) TableName() string {
	return "users" // 或 "members"，取決於你的表名
}
