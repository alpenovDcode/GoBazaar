package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/alpewa/GoBazaar/internal/common/models"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// JWTUtil предоставляет функциональность для работы с JWT токенами
type JWTUtil struct {
	secretKey            []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	issuer               string
}

// NewJWTUtil создает новый экземпляр JWTUtil
func NewJWTUtil(secretKey string, accessTokenHours, refreshTokenDays int, issuer string) *JWTUtil {
	return &JWTUtil{
		secretKey:            []byte(secretKey),
		accessTokenDuration:  time.Duration(accessTokenHours) * time.Hour,
		refreshTokenDuration: time.Duration(refreshTokenDays) * 24 * time.Hour,
		issuer:               issuer,
	}
}

// GenerateAccessToken генерирует access token для пользователя
func (j *JWTUtil) GenerateAccessToken(user *models.User) (string, *models.JWTClaims, error) {
	if user == nil {
		return "", nil, errors.New("user cannot be nil")
	}

	now := time.Now()
	expiresAt := now.Add(j.accessTokenDuration)

	claims := &models.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   fmt.Sprintf("%d", user.ID),
			Audience:  []string{"gobazaar"},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("%d-%d", user.ID, now.Unix()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, claims, nil
}

// GenerateRefreshToken генерирует refresh token
func (j *JWTUtil) GenerateRefreshToken(user *models.User) (string, time.Time, error) {
	if user == nil {
		return "", time.Time{}, errors.New("user cannot be nil")
	}

	now := time.Now()
	expiresAt := now.Add(j.refreshTokenDuration)

	claims := &jwt.RegisteredClaims{
		Issuer:    j.issuer,
		Subject:   fmt.Sprintf("refresh-%d", user.ID),
		Audience:  []string{"gobazaar-refresh"},
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        fmt.Sprintf("ref-%d-%d", user.ID, now.Unix()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// ValidateAccessToken проверяет и парсит access token
func (j *JWTUtil) ValidateAccessToken(tokenString string) (*models.JWTClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token cannot be empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, errors.New("token not valid yet")
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New("token is malformed")
		}
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*models.JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Дополнительные проверки
	if claims.Issuer != j.issuer {
		return nil, errors.New("invalid token issuer")
	}

	if claims.UserID == 0 {
		return nil, errors.New("invalid user ID in token")
	}

	if claims.Email == "" {
		return nil, errors.New("invalid email in token")
	}

	return claims, nil
}

// ValidateRefreshToken проверяет refresh token
func (j *JWTUtil) ValidateRefreshToken(tokenString string) (*jwt.RegisteredClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token cannot be empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("refresh token has expired")
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, errors.New("refresh token not valid yet")
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New("refresh token is malformed")
		}
		return nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token claims")
	}

	// Дополнительные проверки
	if claims.Issuer != j.issuer {
		return nil, errors.New("invalid refresh token issuer")
	}

	// Проверяем, что это refresh token
	if len(claims.Audience) == 0 || claims.Audience[0] != "gobazaar-refresh" {
		return nil, errors.New("invalid refresh token audience")
	}

	return claims, nil
}

// ExtractUserIDFromRefreshToken извлекает ID пользователя из refresh token
func (j *JWTUtil) ExtractUserIDFromRefreshToken(claims *jwt.RegisteredClaims) (uint, error) {
	if claims == nil {
		return 0, errors.New("claims cannot be nil")
	}

	// Subject имеет формат "refresh-{userID}"
	var userID uint
	if _, err := fmt.Sscanf(claims.Subject, "refresh-%d", &userID); err != nil {
		return 0, fmt.Errorf("failed to extract user ID from refresh token: %w", err)
	}

	if userID == 0 {
		return 0, errors.New("invalid user ID in refresh token")
	}

	return userID, nil
}

// GetTokenRemainingTime возвращает оставшееся время жизни токена
func (j *JWTUtil) GetTokenRemainingTime(claims *models.JWTClaims) time.Duration {
	if claims == nil || claims.ExpiresAt == nil {
		return 0
	}

	now := time.Now()
	expiresAt := claims.ExpiresAt.Time

	if expiresAt.Before(now) {
		return 0
	}

	return expiresAt.Sub(now)
}

// IsTokenExpiringSoon проверяет, истекает ли токен скоро
func (j *JWTUtil) IsTokenExpiringSoon(claims *models.JWTClaims, threshold time.Duration) bool {
	remaining := j.GetTokenRemainingTime(claims)
	return remaining <= threshold && remaining > 0
}

// GenerateTokenPair генерирует пару access и refresh токенов
func (j *JWTUtil) GenerateTokenPair(user *models.User) (string, string, time.Time, error) {
	// Генерируем access token
	accessToken, _, err := j.GenerateAccessToken(user)
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Генерируем refresh token
	refreshToken, expiresAt, err := j.GenerateRefreshToken(user)
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, expiresAt, nil
}

// GetAccessTokenDuration возвращает длительность жизни access token
func (j *JWTUtil) GetAccessTokenDuration() time.Duration {
	return j.accessTokenDuration
}

// GetRefreshTokenDuration возвращает длительность жизни refresh token
func (j *JWTUtil) GetRefreshTokenDuration() time.Duration {
	return j.refreshTokenDuration
}

// GetIssuer возвращает issuer токенов
func (j *JWTUtil) GetIssuer() string {
	return j.issuer
}
