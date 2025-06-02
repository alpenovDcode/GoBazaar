package models

import (
	"testing"
	"time"
)

func TestUser_TableName(t *testing.T) {
	user := User{}
	expected := "users"

	if got := user.TableName(); got != expected {
		t.Errorf("TableName() = %v, want %v", got, expected)
	}
}

func TestUser_Structure(t *testing.T) {
	user := User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if user.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", user.ID)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected Email to be 'test@example.com', got %s", user.Email)
	}

	if user.PasswordHash != "hashed_password" {
		t.Errorf("Expected PasswordHash to be 'hashed_password', got %s", user.PasswordHash)
	}
}
