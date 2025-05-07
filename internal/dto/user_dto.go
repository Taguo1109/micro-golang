package dto

/**
 * @File: user_dto.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/5/7 下午2:41
 * @Software: GoLand
 * @Version:  1.0
 */

type UserDTO struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	Role         string `json:"role"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserLoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserLoginResponseDTO struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type UserLogoutDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserRegisterDTO struct {
	Email    string `json:"email" binding:"required,email" validateMsg:"required=Email 為必填,email=Email 格式錯誤" example:"test@example.com"`
	Username string `json:"username" binding:"required,username_validation" validateMsg:"required=使用者名稱為必填,username_validation=使用者名稱只能是英文與數字，且長度為 6~20 字" example:"testUser01"`
	Password string `json:"password" binding:"required,pwd_validation" validateMsg:"required=密碼為必填,pwd_validation=密碼需包含至少一個大寫與一個小寫字母，且長度 6~30 字" example:"P@ssw0rd"`
	Role     string `json:"role" binding:"required,oneof=User Admin SuperAdmin" validateMsg:"required=角色為必填,oneof=角色只能是 User、Admin 或 SuperAdmin" example:"User"`
}
