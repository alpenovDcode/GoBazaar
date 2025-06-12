package utils

import (
	"testing"
	"time"

	"github.com/alpewa/GoBazaar/internal/common/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestUser creates a test user
func createTestUser() *models.User {
	return &models.User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      models.RoleCustomer,
		IsActive:  true,
	}
}

func TestJWTUtil_GenerateAccessToken(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key-that-is-long-enough", 1, 7, "test-issuer")
	user := createTestUser()

	t.Run("successful generation", func(t *testing.T) {
		token, claims, err := jwtUtil.GenerateAccessToken(user)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotNil(t, claims)
		assert.Equal(t, user.ID, claims.UserID)
		assert.Equal(t, user.Email, claims.Email)
		assert.Equal(t, user.Role, claims.Role)
		assert.Equal(t, "test-issuer", claims.Issuer)
	})

	t.Run("nil user", func(t *testing.T) {
		_, _, err := jwtUtil.GenerateAccessToken(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be nil")
	})
}

func TestJWTUtil_GenerateRefreshToken(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key-that-is-long-enough", 1, 7, "test-issuer")
	user := createTestUser()

	t.Run("successful generation", func(t *testing.T) {
		token, expiresAt, err := jwtUtil.GenerateRefreshToken(user)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.True(t, expiresAt.After(time.Now()))

		// Check that token expires approximately in 7 days
		expectedExpiry := time.Now().Add(7 * 24 * time.Hour)
		assert.WithinDuration(t, expectedExpiry, expiresAt, time.Hour)
	})

	t.Run("nil user", func(t *testing.T) {
		_, _, err := jwtUtil.GenerateRefreshToken(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be nil")
	})
}

func TestJWTUtil_ValidateAccessToken(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key-that-is-long-enough", 1, 7, "test-issuer")
	user := createTestUser()

	t.Run("valid token", func(t *testing.T) {
		token, _, err := jwtUtil.GenerateAccessToken(user)
		require.NoError(t, err)

		claims, err := jwtUtil.ValidateAccessToken(token)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, claims.UserID)
		assert.Equal(t, user.Email, claims.Email)
		assert.Equal(t, user.Role, claims.Role)
	})

	t.Run("empty token", func(t *testing.T) {
		_, err := jwtUtil.ValidateAccessToken("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := jwtUtil.ValidateAccessToken("invalid.token.here")
		assert.Error(t, err)
	})

	t.Run("token with wrong secret", func(t *testing.T) {
		wrongJwtUtil := NewJWTUtil("wrong-secret-key-that-is-long-enough", 1, 7, "test-issuer")
		token, _, err := wrongJwtUtil.GenerateAccessToken(user)
		require.NoError(t, err)

		_, err = jwtUtil.ValidateAccessToken(token)
		assert.Error(t, err)
	})

	t.Run("token with wrong issuer", func(t *testing.T) {
		wrongJwtUtil := NewJWTUtil("test-secret-key-that-is-long-enough", 1, 7, "wrong-issuer")
		token, _, err := wrongJwtUtil.GenerateAccessToken(user)
		require.NoError(t, err)

		_, err = jwtUtil.ValidateAccessToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token issuer")
	})
}

func TestJWTUtil_ValidateRefreshToken(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key-that-is-long-enough", 1, 7, "test-issuer")
	user := createTestUser()

	t.Run("valid refresh token", func(t *testing.T) {
		token, _, err := jwtUtil.GenerateRefreshToken(user)
		require.NoError(t, err)

		claims, err := jwtUtil.ValidateRefreshToken(token)
		assert.NoError(t, err)
		assert.Equal(t, "test-issuer", claims.Issuer)
		assert.Contains(t, claims.Audience, "gobazaar-refresh")
	})

	t.Run("empty token", func(t *testing.T) {
		_, err := jwtUtil.ValidateRefreshToken("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("access token as refresh token", func(t *testing.T) {
		token, _, err := jwtUtil.GenerateAccessToken(user)
		require.NoError(t, err)

		_, err = jwtUtil.ValidateRefreshToken(token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid refresh token audience")
	})
}

func TestJWTUtil_ExtractUserIDFromRefreshToken(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key-that-is-long-enough", 1, 7, "test-issuer")
	user := createTestUser()

	t.Run("valid extraction", func(t *testing.T) {
		token, _, err := jwtUtil.GenerateRefreshToken(user)
		require.NoError(t, err)

		claims, err := jwtUtil.ValidateRefreshToken(token)
		require.NoError(t, err)

		userID, err := jwtUtil.ExtractUserIDFromRefreshToken(claims)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, userID)
	})

	t.Run("nil claims", func(t *testing.T) {
		_, err := jwtUtil.ExtractUserIDFromRefreshToken(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be nil")
	})
}

func TestJWTUtil_GetTokenRemainingTime(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key-that-is-long-enough", 1, 7, "test-issuer")
	user := createTestUser()

	t.Run("valid token", func(t *testing.T) {
		_, claims, err := jwtUtil.GenerateAccessToken(user)
		require.NoError(t, err)

		remaining := jwtUtil.GetTokenRemainingTime(claims)
		assert.True(t, remaining > 0)
		assert.True(t, remaining <= time.Hour) // Token created for 1 hour
	})

	t.Run("nil claims", func(t *testing.T) {
		remaining := jwtUtil.GetTokenRemainingTime(nil)
		assert.Equal(t, time.Duration(0), remaining)
	})
}

func TestJWTUtil_IsTokenExpiringSoon(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key-that-is-long-enough", 1, 7, "test-issuer")
	user := createTestUser()

	t.Run("not expiring soon", func(t *testing.T) {
		_, claims, err := jwtUtil.GenerateAccessToken(user)
		require.NoError(t, err)

		isExpiringSoon := jwtUtil.IsTokenExpiringSoon(claims, 10*time.Minute)
		assert.False(t, isExpiringSoon) // Token lives for an hour, checking for 10 minutes
	})

	t.Run("expiring soon", func(t *testing.T) {
		_, claims, err := jwtUtil.GenerateAccessToken(user)
		require.NoError(t, err)

		isExpiringSoon := jwtUtil.IsTokenExpiringSoon(claims, 2*time.Hour)
		assert.True(t, isExpiringSoon) // Token lives for an hour, checking for 2 hours
	})
}

func TestJWTUtil_GenerateTokenPair(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key-that-is-long-enough", 1, 7, "test-issuer")
	user := createTestUser()

	t.Run("successful generation", func(t *testing.T) {
		accessToken, refreshToken, expiresAt, err := jwtUtil.GenerateTokenPair(user)
		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)
		assert.True(t, expiresAt.After(time.Now()))

		// Check that tokens are valid
		_, err = jwtUtil.ValidateAccessToken(accessToken)
		assert.NoError(t, err)

		_, err = jwtUtil.ValidateRefreshToken(refreshToken)
		assert.NoError(t, err)
	})

	t.Run("nil user", func(t *testing.T) {
		_, _, _, err := jwtUtil.GenerateTokenPair(nil)
		assert.Error(t, err)
	})
}

func TestJWTUtil_GettersAndSetters(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key-that-is-long-enough", 2, 14, "test-issuer")

	t.Run("get durations", func(t *testing.T) {
		assert.Equal(t, 2*time.Hour, jwtUtil.GetAccessTokenDuration())
		assert.Equal(t, 14*24*time.Hour, jwtUtil.GetRefreshTokenDuration())
		assert.Equal(t, "test-issuer", jwtUtil.GetIssuer())
	})
}
