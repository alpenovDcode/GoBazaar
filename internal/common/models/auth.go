package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// UserRole defines user roles
type UserRole string

const (
	RoleCustomer  UserRole = "customer"
	RoleAdmin     UserRole = "admin"
	RoleModerator UserRole = "moderator"
)

// User represents a user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Password  string         `json:"-" gorm:"not null" validate:"required,min=8"`
	FirstName string         `json:"first_name" gorm:"not null" validate:"required,min=2,max=50"`
	LastName  string         `json:"last_name" gorm:"not null" validate:"required,min=2,max=50"`
	Role      UserRole       `json:"role" gorm:"default:'customer'" validate:"required"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserResponse is used for API responses (without password)
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      UserRole  `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required,min=8"`
	FirstName string   `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string   `json:"last_name" validate:"required,min=2,max=50"`
	Role      UserRole `json:"role,omitempty"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    time.Time    `json:"expires_at"`
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID uint     `json:"user_id"`
	Email  string   `json:"email"`
	Role   UserRole `json:"role"`
	jwt.RegisteredClaims
}

// RefreshToken represents refresh token in database
type RefreshToken struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	Token     string         `json:"token" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null"`
	IsRevoked bool           `json:"is_revoked" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	User      User           `json:"-" gorm:"foreignKey:UserID"`
}

// ChangePasswordRequest represents change password request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// ResetPasswordRequest represents reset password request
type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// PasswordResetToken represents password reset token
type PasswordResetToken struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	Token     string         `json:"token" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null"`
	IsUsed    bool           `json:"is_used" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	User      User           `json:"-" gorm:"foreignKey:UserID"`
}
