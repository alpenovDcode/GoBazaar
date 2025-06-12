package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/alpewa/GoBazaar/internal/auth/config"
	"github.com/alpewa/GoBazaar/internal/auth/repository"
	"github.com/alpewa/GoBazaar/internal/auth/utils"
	"github.com/alpewa/GoBazaar/internal/common/models"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotActive      = errors.New("user account is not active")
	ErrValidationFailed   = errors.New("validation failed")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token has expired")
)

// AuthService интерфейс для auth сервиса
type AuthService interface {
	// Аутентификация
	Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshTokenString string) (*models.AuthResponse, error)
	Logout(ctx context.Context, refreshTokenString string) error
	ValidateAccessToken(tokenString string) (*models.JWTClaims, error)

	// Управление пользователями
	GetUser(ctx context.Context, userID uint) (*models.UserResponse, error)
	UpdateUser(ctx context.Context, userID uint, updates map[string]interface{}) (*models.UserResponse, error)
	ChangePassword(ctx context.Context, userID uint, req *models.ChangePasswordRequest) error

	// Администрирование
	GetAllUsers(ctx context.Context, limit, offset int) ([]*models.UserResponse, int64, error)
	SearchUsers(ctx context.Context, query string, limit, offset int) ([]*models.UserResponse, int64, error)
	DeactivateUser(ctx context.Context, userID uint) error
	ActivateUser(ctx context.Context, userID uint) error
}

// authService реализация auth сервиса
type authService struct {
	userRepo  repository.UserRepository
	jwtUtil   *utils.JWTUtil
	validator *validator.Validate
	config    *config.Config
}

// NewAuthService создает новый экземпляр auth сервиса
func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	jwtUtil := utils.NewJWTUtil(
		cfg.JWTSecret,
		cfg.JWTExpirationHours,
		cfg.RefreshTokenExpDays,
		"gobazaar-auth",
	)

	return &authService{
		userRepo:  userRepo,
		jwtUtil:   jwtUtil,
		validator: validator.New(),
		config:    cfg,
	}
}

// Register регистрирует нового пользователя
func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Валидация запроса
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidationFailed, err)
	}

	// Проверяем минимальную длину пароля
	if len(req.Password) < s.config.PasswordMinLength {
		return nil, fmt.Errorf("password must be at least %d characters long", s.config.PasswordMinLength)
	}

	// Проверяем существование пользователя
	exists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, ErrUserExists
	}

	// Хешируем пароль
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Устанавливаем роль по умолчанию
	role := req.Role
	if role == "" {
		role = models.RoleCustomer
	}

	// Создаем пользователя
	user := &models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      role,
		IsActive:  true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Генерируем токены
	return s.generateAuthResponse(user)
}

// Login аутентифицирует пользователя
func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	// Валидация запроса
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidationFailed, err)
	}

	// Получаем пользователя по email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, errors.New("user not found")) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Проверяем активность пользователя
	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Проверяем пароль
	if err := s.checkPassword(req.Password, user.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Обновляем время последнего входа
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Логируем ошибку, но не прерываем процесс
		fmt.Printf("Failed to update last login for user %d: %v\n", user.ID, err)
	}

	// Генерируем токены
	return s.generateAuthResponse(user)
}

// RefreshToken обновляет access token используя refresh token
func (s *authService) RefreshToken(ctx context.Context, refreshTokenString string) (*models.AuthResponse, error) {
	// Валидируем refresh token
	refreshClaims, err := s.jwtUtil.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Извлекаем ID пользователя
	userID, err := s.jwtUtil.ExtractUserIDFromRefreshToken(refreshClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to extract user ID: %w", err)
	}

	// Получаем пользователя
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Проверяем активность пользователя
	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Генерируем новые токены
	return s.generateAuthResponse(user)
}

// Logout отзывает токены пользователя (заглушка)
func (s *authService) Logout(ctx context.Context, refreshTokenString string) error {
	// В данной реализации используем stateless JWT
	// В production можно добавить blacklist для токенов
	return nil
}

// ValidateAccessToken валидирует access token
func (s *authService) ValidateAccessToken(tokenString string) (*models.JWTClaims, error) {
	return s.jwtUtil.ValidateAccessToken(tokenString)
}

// GetUser получает информацию о пользователе
func (s *authService) GetUser(ctx context.Context, userID uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

// UpdateUser обновляет информацию о пользователе
func (s *authService) UpdateUser(ctx context.Context, userID uint, updates map[string]interface{}) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Обновляем разрешенные поля
	if firstName, ok := updates["first_name"].(string); ok && firstName != "" {
		user.FirstName = firstName
	}
	if lastName, ok := updates["last_name"].(string); ok && lastName != "" {
		user.LastName = lastName
	}
	if email, ok := updates["email"].(string); ok && email != "" {
		user.Email = email
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

// ChangePassword изменяет пароль пользователя
func (s *authService) ChangePassword(ctx context.Context, userID uint, req *models.ChangePasswordRequest) error {
	// Валидация запроса
	if err := s.validator.Struct(req); err != nil {
		return fmt.Errorf("%w: %v", ErrValidationFailed, err)
	}

	// Проверяем длину нового пароля
	if len(req.NewPassword) < s.config.PasswordMinLength {
		return fmt.Errorf("new password must be at least %d characters long", s.config.PasswordMinLength)
	}

	// Получаем пользователя
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Проверяем текущий пароль
	if err := s.checkPassword(req.CurrentPassword, user.Password); err != nil {
		return ErrInvalidCredentials
	}

	// Хешируем новый пароль
	hashedPassword, err := s.hashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Обновляем пароль
	user.Password = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// GetAllUsers получает список всех пользователей (для администраторов)
func (s *authService) GetAllUsers(ctx context.Context, limit, offset int) ([]*models.UserResponse, int64, error) {
	users, total, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	responses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		response := user.ToResponse()
		responses[i] = &response
	}

	return responses, total, nil
}

// SearchUsers ищет пользователей по имени
func (s *authService) SearchUsers(ctx context.Context, query string, limit, offset int) ([]*models.UserResponse, int64, error) {
	users, total, err := s.userRepo.SearchByName(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}

	responses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		response := user.ToResponse()
		responses[i] = &response
	}

	return responses, total, nil
}

// DeactivateUser деактивирует пользователя
func (s *authService) DeactivateUser(ctx context.Context, userID uint) error {
	if err := s.userRepo.DeactivateUser(ctx, userID); err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}
	return nil
}

// ActivateUser активирует пользователя
func (s *authService) ActivateUser(ctx context.Context, userID uint) error {
	if err := s.userRepo.ActivateUser(ctx, userID); err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}
	return nil
}

// generateAuthResponse генерирует ответ с токенами
func (s *authService) generateAuthResponse(user *models.User) (*models.AuthResponse, error) {
	// Генерируем access token
	accessToken, accessClaims, err := s.jwtUtil.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Генерируем refresh token
	refreshToken, refreshExpiresAt, err := s.jwtUtil.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Создаем refresh token в базе данных (если нужно)
	// В данной реализации используем stateless JWT

	response := &models.AuthResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessClaims.ExpiresAt.Time,
	}

	// Можно добавить информацию о refresh token
	_ = refreshExpiresAt

	return response, nil
}

// hashPassword хеширует пароль
func (s *authService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.config.BCryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// checkPassword проверяет пароль
func (s *authService) checkPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
