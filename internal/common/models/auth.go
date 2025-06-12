package models

import (
	"time"

	"gorm.io/gorm"
)

// LoginAttempt represents a login attempt (for rate limiting purposes)
type LoginAttempt struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"not null;index"`
	IPAddress string         `json:"ip_address" gorm:"not null"`
	UserAgent string         `json:"user_agent"`
	Success   bool           `json:"success" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Session represents a user session
type Session struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"not null;index"`
	SessionID  string         `json:"session_id" gorm:"uniqueIndex;not null"`
	IPAddress  string         `json:"ip_address"`
	UserAgent  string         `json:"user_agent"`
	ExpiresAt  time.Time      `json:"expires_at" gorm:"not null"`
	IsActive   bool           `json:"is_active" gorm:"default:true"`
	LastUsedAt time.Time      `json:"last_used_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	User       User           `json:"-" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for LoginAttempt
func (LoginAttempt) TableName() string {
	return "login_attempts"
}

// TableName returns the table name for Session
func (Session) TableName() string {
	return "sessions"
}
