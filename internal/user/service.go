package user

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"micro-golang/internal/config"
	"micro-golang/internal/dto"
	"micro-golang/internal/models"
	"strings"
	"time"
)

/**
 * @File: service.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/23 上午11:12
 * @Software: GoLand
 * @Version:  1.0
 */

// Service 負責處理用戶相關的業務邏輯
type Service struct {
	db  *gorm.DB
	rdb *redis.Client // 假設你的 config.RDB 是 *redis.Client 的類型，或者你可以直接用 config.RDB
}

// NewService 創建 Service 實例
func NewService(db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{
		db:  db,
		rdb: rdb, // 或者直接在方法中使用 config.RDB
	}
}

// Custom error types for service layer
var (
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailInUse       = errors.New("email is already in use")
	ErrUpdateNoChanges  = errors.New("no fields provided for update or no changes detected")
	ErrValidationFailed = errors.New("validation failed") // For business rule validation
)

// UpdateUserProfile 更新用戶資料 (Username, Email)
func (s *Service) UpdateUserProfile(ctx context.Context, currentEmail string, req dto.UserUpdateProfileDTO) (*dto.UserLoginResponseDTO, error) {
	if req.Username == nil && req.Email == nil {
		return nil, ErrUpdateNoChanges
	}

	var user models.User
	if err := s.db.WithContext(ctx).Where("email = ?", currentEmail).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err // Return generic DB error
	}

	oldCacheKey := "user:" + user.Email
	emailChanged := false
	updates := make(map[string]interface{})

	if req.Username != nil && *req.Username != user.Username {
		// TODO: 在此處添加更詳細的 Username 業務驗證邏輯 (如果需要)
		// 例如: if len(*req.Username) < 3 { return nil, fmt.Errorf("%w: username too short", ErrValidationFailed) }
		updates["username"] = *req.Username
		user.Username = *req.Username
	}

	if req.Email != nil && strings.ToLower(*req.Email) != strings.ToLower(user.Email) {
		newEmail := strings.ToLower(*req.Email)
		// TODO: 在此處添加更詳細的 Email 業務驗證邏輯 (如果需要)

		var existingUserWithNewEmail models.User
		err := s.db.WithContext(ctx).Where("email = ? AND id != ?", newEmail, user.ID).First(&existingUserWithNewEmail).Error
		if err == nil {
			return nil, ErrEmailInUse
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err // DB 查詢本身出錯
		}

		updates["email"] = newEmail
		user.Email = newEmail
		emailChanged = true
	}

	if len(updates) == 0 {
		// 如果 DTO 有值但值與現有資料相同，這裡也會被觸發
		// 可以回傳更新後的 user DTO，並標註 "no effective changes"
		// 或者嚴格一點，若 DTO 有提供值但值未變，也視為 ErrUpdateNoChanges
		// 這裡我們假設只要有欄位被提供，就算值相同也繼續，讓 DB 的 updated_at 更新
		// 若要嚴格 "no changes"，則需調整此處邏輯
		// return nil, ErrUpdateNoChanges // 如果 DTO 提供的值與現有值相同，則視為無變更
	}
	// 即使沒有欄位需要更新，但若 DTO 有提供，也可能只是值相同，
	// 為了確保 updated_at 被更新 (如果這是期望行為)，我們可以繼續執行。
	// 如果嚴格要求只有在值不同的情況下才更新，那麼上面的 len(updates) == 0 判斷應該更早。
	// 目前的邏輯是，只要 updates map 中有值，就會執行更新。
	// 若 DTO 中有值，但與 DB 中相同，updates map 中就不會有該欄位。

	// 如果真的沒有任何欄位被賦予新值 (DTO 有值但與 DB 相同，或 DTO 欄位為 nil)
	if len(updates) == 0 {
		// 返回當前用戶信息，表示沒有實際更改
		safeUser := dto.UserLoginResponseDTO{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Role:     user.Role,
		}
		return &safeUser, ErrUpdateNoChanges // 使用一個特定的 error 或 nil 來表示無變更但操作成功
	}

	updates["updated_at"] = time.Now() // GORM 通常會自動處理，但顯式指定也無妨

	if err := s.db.WithContext(ctx).Model(&user).Updates(updates).Error; err != nil {
		return nil, err // Return generic DB error
	}

	// 更新快取邏輯
	if emailChanged {
		_, errCacheDel := s.rdb.Del(config.Ctx, oldCacheKey).Result() // 使用注入的 rdb 或全局的 config.Ctx
		if errCacheDel != nil {
			// Log cache deletion error, but proceed as DB update was successful
			// log.Printf("Warning: Failed to delete old cache key %s: %v", oldCacheKey, errCacheDel)
		}
	}

	newCacheKey := "user:" + user.Email
	updatedSafeUserDTO := dto.UserLoginResponseDTO{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Role:     user.Role,
	}
	userBytes, err := json.Marshal(updatedSafeUserDTO)
	if err != nil {
		// Log marshalling error, proceed but indicate cache issue
		// log.Printf("Warning: Failed to marshal user data for cache: %v", err)
		return &updatedSafeUserDTO, errors.New("user updated successfully, but cache marshal failed")
	}

	err = s.rdb.Set(config.Ctx, newCacheKey, userBytes, 10*time.Minute).Err()
	if err != nil {
		// Log cache set error
		// log.Printf("Warning: Failed to set user cache for %s: %v", newCacheKey, err)
		return &updatedSafeUserDTO, errors.New("user updated successfully, but cache set failed")
	}

	return &updatedSafeUserDTO, nil
}
