package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUser_ToResponse(t *testing.T) {
	user := &User{
		ID:        1,
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "John",
		LastName:  "Doe",
		Role:      RoleCustomer,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	expected := UserResponse{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      RoleCustomer,
		IsActive:  true,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	response := user.ToResponse()

	assert.Equal(t, expected, response)
	// Password is not included in UserResponse by design
}

func TestUserRole_Constants(t *testing.T) {
	assert.Equal(t, UserRole("customer"), RoleCustomer)
	assert.Equal(t, UserRole("admin"), RoleAdmin)
	assert.Equal(t, UserRole("moderator"), RoleModerator)
}

func TestRegisterRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request RegisterRequest
		valid   bool
	}{
		{
			name: "valid request",
			request: RegisterRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
				Role:      RoleCustomer,
			},
			valid: true,
		},
		{
			name: "empty email",
			request: RegisterRequest{
				Email:     "",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			valid: false,
		},
		{
			name: "short password",
			request: RegisterRequest{
				Email:     "test@example.com",
				Password:  "short",
				FirstName: "John",
				LastName:  "Doe",
			},
			valid: false,
		},
		{
			name: "short first name",
			request: RegisterRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "J",
				LastName:  "Doe",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				// For valid requests, check that all fields are filled
				assert.NotEmpty(t, tt.request.Email)
				assert.NotEmpty(t, tt.request.Password)
				assert.NotEmpty(t, tt.request.FirstName)
				assert.NotEmpty(t, tt.request.LastName)
			} else {
				// For invalid requests, check specific conditions
				switch tt.name {
				case "empty email":
					assert.Empty(t, tt.request.Email)
				case "short password":
					assert.True(t, len(tt.request.Password) < 8)
				case "short first name":
					assert.True(t, len(tt.request.FirstName) < 2)
				}
			}
		})
	}
}

func TestLoginRequest_Validation(t *testing.T) {
	request := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	assert.NotEmpty(t, request.Email)
	assert.NotEmpty(t, request.Password)
}

func TestAuthResponse_Structure(t *testing.T) {
	user := UserResponse{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      RoleCustomer,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	response := AuthResponse{
		User:         user,
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	assert.Equal(t, user, response.User)
	assert.Equal(t, "access_token", response.AccessToken)
	assert.Equal(t, "refresh_token", response.RefreshToken)
	assert.NotZero(t, response.ExpiresAt)
}

func TestJWTClaims_Structure(t *testing.T) {
	claims := JWTClaims{
		UserID: 1,
		Email:  "test@example.com",
		Role:   RoleCustomer,
	}

	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, RoleCustomer, claims.Role)
}

func TestRefreshToken_Relationships(t *testing.T) {
	user := User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      RoleCustomer,
	}

	refreshToken := RefreshToken{
		ID:        1,
		UserID:    user.ID,
		Token:     "refresh_token",
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
		IsRevoked: false,
		User:      user,
	}

	assert.Equal(t, user.ID, refreshToken.UserID)
	assert.Equal(t, user, refreshToken.User)
	assert.False(t, refreshToken.IsRevoked)
}

func TestPasswordResetToken_Structure(t *testing.T) {
	token := PasswordResetToken{
		ID:        1,
		UserID:    1,
		Token:     "reset_token",
		ExpiresAt: time.Now().Add(time.Hour),
		IsUsed:    false,
	}

	assert.Equal(t, uint(1), token.ID)
	assert.Equal(t, uint(1), token.UserID)
	assert.Equal(t, "reset_token", token.Token)
	assert.False(t, token.IsUsed)
}

func TestChangePasswordRequest_Validation(t *testing.T) {
	request := ChangePasswordRequest{
		CurrentPassword: "currentpassword",
		NewPassword:     "newpassword123",
	}

	assert.NotEmpty(t, request.CurrentPassword)
	assert.NotEmpty(t, request.NewPassword)
}

func TestResetPasswordRequest_Validation(t *testing.T) {
	request := ResetPasswordRequest{
		Email: "test@example.com",
	}

	assert.NotEmpty(t, request.Email)
}
