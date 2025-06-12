package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alpewa/GoBazaar/internal/common/models"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUserExists    = errors.New("user already exists")
	ErrTokenNotFound = errors.New("refresh token not found")
	ErrTokenExpired  = errors.New("refresh token expired")
	ErrTokenRevoked  = errors.New("refresh token revoked")
)

// UserRepository interface for working with users
type UserRepository interface {
	// CRUD operations
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]*models.User, int64, error)

	// Search and filtering
	SearchByName(ctx context.Context, name string, limit, offset int) ([]*models.User, int64, error)
	GetByRole(ctx context.Context, role models.UserRole, limit, offset int) ([]*models.User, int64, error)
	GetActiveUsers(ctx context.Context, limit, offset int) ([]*models.User, int64, error)

	// Authentication
	UpdateLastLogin(ctx context.Context, userID uint) error
	DeactivateUser(ctx context.Context, userID uint) error
	ActivateUser(ctx context.Context, userID uint) error

	// Statistics
	CountUsers(ctx context.Context) (int64, error)
	CountUsersByRole(ctx context.Context, role models.UserRole) (int64, error)
	CountActiveUsers(ctx context.Context) (int64, error)

	// Uniqueness checks
	EmailExists(ctx context.Context, email string) (bool, error)
	EmailExistsForOtherUser(ctx context.Context, email string, userID uint) (bool, error)
}

// userRepository реализация репозитория пользователей
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository создает новый экземпляр репозитория пользователей
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create создает нового пользователя
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	// Проверяем уникальность email
	exists, err := r.EmailExists(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("failed to check email uniqueness: %w", err)
	}
	if exists {
		return errors.New("user with this email already exists")
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID получает пользователя по ID
func (r *userRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	if id == 0 {
		return nil, errors.New("invalid user ID")
	}

	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

// GetByEmail получает пользователя по email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// Update обновляет пользователя
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	if user.ID == 0 {
		return errors.New("user ID cannot be zero")
	}

	// Проверяем уникальность email для других пользователей
	exists, err := r.EmailExistsForOtherUser(ctx, user.Email, user.ID)
	if err != nil {
		return fmt.Errorf("failed to check email uniqueness: %w", err)
	}
	if exists {
		return errors.New("user with this email already exists")
	}

	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete удаляет пользователя (soft delete)
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid user ID")
	}

	if err := r.db.WithContext(ctx).Delete(&models.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// List получает список пользователей с пагинацией
func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	// Получаем общее количество
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Получаем пользователей с пагинацией
	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

// SearchByName ищет пользователей по имени
func (r *userRepository) SearchByName(ctx context.Context, name string, limit, offset int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	searchPattern := "%" + name + "%"
	query := r.db.WithContext(ctx).Model(&models.User{}).Where(
		"first_name ILIKE ? OR last_name ILIKE ?",
		searchPattern, searchPattern,
	)

	// Получаем общее количество
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users by name: %w", err)
	}

	// Получаем пользователей с пагинацией
	if err := query.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search users by name: %w", err)
	}

	return users, total, nil
}

// GetByRole получает пользователей по роли
func (r *userRepository) GetByRole(ctx context.Context, role models.UserRole, limit, offset int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	query := r.db.WithContext(ctx).Model(&models.User{}).Where("role = ?", role)

	// Получаем общее количество
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users by role: %w", err)
	}

	// Получаем пользователей с пагинацией
	if err := query.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get users by role: %w", err)
	}

	return users, total, nil
}

// GetActiveUsers получает активных пользователей
func (r *userRepository) GetActiveUsers(ctx context.Context, limit, offset int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	query := r.db.WithContext(ctx).Model(&models.User{}).Where("is_active = ?", true)

	// Получаем общее количество
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count active users: %w", err)
	}

	// Получаем пользователей с пагинацией
	if err := query.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get active users: %w", err)
	}

	return users, total, nil
}

// UpdateLastLogin обновляет время последнего входа
func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uint) error {
	if userID == 0 {
		return errors.New("invalid user ID")
	}

	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("updated_at", time.Now()).Error; err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// DeactivateUser деактивирует пользователя
func (r *userRepository) DeactivateUser(ctx context.Context, userID uint) error {
	if userID == 0 {
		return errors.New("invalid user ID")
	}

	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("is_active", false).Error; err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	return nil
}

// ActivateUser активирует пользователя
func (r *userRepository) ActivateUser(ctx context.Context, userID uint) error {
	if userID == 0 {
		return errors.New("invalid user ID")
	}

	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("is_active", true).Error; err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}

	return nil
}

// CountUsers подсчитывает общее количество пользователей
func (r *userRepository) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// CountUsersByRole подсчитывает пользователей по роли
func (r *userRepository) CountUsersByRole(ctx context.Context, role models.UserRole) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("role = ?", role).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users by role: %w", err)
	}
	return count, nil
}

// CountActiveUsers подсчитывает активных пользователей
func (r *userRepository) CountActiveUsers(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("is_active = ?", true).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count active users: %w", err)
	}
	return count, nil
}

// EmailExists проверяет существование email
func (r *userRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	if email == "" {
		return false, errors.New("email cannot be empty")
	}

	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}

// EmailExistsForOtherUser проверяет существование email у других пользователей
func (r *userRepository) EmailExistsForOtherUser(ctx context.Context, email string, userID uint) (bool, error) {
	if email == "" {
		return false, errors.New("email cannot be empty")
	}
	if userID == 0 {
		return false, errors.New("invalid user ID")
	}

	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ? AND id != ?", email, userID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check email existence for other users: %w", err)
	}

	return count > 0, nil
}

// CreateRefreshToken creates a new refresh token
func (r *userRepository) CreateRefreshToken(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

// GetRefreshToken gets a refresh token by token string
func (r *userRepository) GetRefreshToken(tokenString string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	if err := r.db.Where("token = ?", tokenString).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}

	// Check if token is expired
	if token.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	// Check if token is revoked
	if token.IsRevoked {
		return nil, ErrTokenRevoked
	}

	return &token, nil
}

// RevokeRefreshToken revokes a refresh token
func (r *userRepository) RevokeRefreshToken(tokenString string) error {
	return r.db.Model(&models.RefreshToken{}).
		Where("token = ?", tokenString).
		Update("is_revoked", true).Error
}

// RevokeAllUserTokens revokes all refresh tokens for a user
func (r *userRepository) RevokeAllUserTokens(userID uint) error {
	return r.db.Model(&models.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("is_revoked", true).Error
}

// CleanupExpiredTokens removes expired refresh tokens
func (r *userRepository) CleanupExpiredTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{}).Error
}
