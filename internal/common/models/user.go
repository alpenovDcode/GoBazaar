package models

import (
	"time"
)

// User представляет пользователя в системе
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName возвращает имя таблицы для модели User
func (User) TableName() string {
	return "users"
}
