package repository

import (
	"context"
	"testing"

	"github.com/alpewa/GoBazaar/internal/common/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB creates an in-memory test database
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Run migrations
	err = db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.PasswordResetToken{},
	)
	require.NoError(t, err)

	return db
}

// createTestUser creates a test user
func createTestUser() *models.User {
	return &models.User{
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "John",
		LastName:  "Doe",
		Role:      models.RoleCustomer,
		IsActive:  true,
	}
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		user := createTestUser()
		err := repo.Create(ctx, user)
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)
		assert.NotZero(t, user.CreatedAt)
	})

	t.Run("duplicate email", func(t *testing.T) {
		user1 := createTestUser()
		user2 := createTestUser()

		err := repo.Create(ctx, user1)
		assert.NoError(t, err)

		err = repo.Create(ctx, user2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("nil user", func(t *testing.T) {
		err := repo.Create(ctx, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be nil")
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("existing user", func(t *testing.T) {
		user := createTestUser()
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		found, err := repo.GetByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user.Email, found.Email)
		assert.Equal(t, user.FirstName, found.FirstName)
	})

	t.Run("non-existing user", func(t *testing.T) {
		_, err := repo.GetByID(ctx, 999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("invalid ID", func(t *testing.T) {
		_, err := repo.GetByID(ctx, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("existing user", func(t *testing.T) {
		user := createTestUser()
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		found, err := repo.GetByEmail(ctx, user.Email)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, found.ID)
		assert.Equal(t, user.FirstName, found.FirstName)
	})

	t.Run("non-existing user", func(t *testing.T) {
		_, err := repo.GetByEmail(ctx, "nonexistent@example.com")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("empty email", func(t *testing.T) {
		_, err := repo.GetByEmail(ctx, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		user := createTestUser()
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		user.FirstName = "Jane"
		user.LastName = "Smith"
		err = repo.Update(ctx, user)
		assert.NoError(t, err)

		updated, err := repo.GetByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Jane", updated.FirstName)
		assert.Equal(t, "Smith", updated.LastName)
	})

	t.Run("nil user", func(t *testing.T) {
		err := repo.Update(ctx, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be nil")
	})

	t.Run("zero ID", func(t *testing.T) {
		user := createTestUser()
		user.ID = 0
		err := repo.Update(ctx, user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be zero")
	})
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("successful delete", func(t *testing.T) {
		user := createTestUser()
		err := repo.Create(ctx, user)
		require.NoError(t, err)

		err = repo.Delete(ctx, user.ID)
		assert.NoError(t, err)

		// Check that user is deleted (soft delete)
		_, err = repo.GetByID(ctx, user.ID)
		assert.Error(t, err)
	})

	t.Run("invalid ID", func(t *testing.T) {
		err := repo.Delete(ctx, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})
}

func TestUserRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create test users
	users := []*models.User{
		{Email: "user1@example.com", Password: "password", FirstName: "User", LastName: "One", Role: models.RoleCustomer, IsActive: true},
		{Email: "user2@example.com", Password: "password", FirstName: "User", LastName: "Two", Role: models.RoleCustomer, IsActive: true},
		{Email: "user3@example.com", Password: "password", FirstName: "User", LastName: "Three", Role: models.RoleAdmin, IsActive: true},
	}

	for _, user := range users {
		err := repo.Create(ctx, user)
		require.NoError(t, err)
	}

	t.Run("list with pagination", func(t *testing.T) {
		result, total, err := repo.List(ctx, 2, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, result, 2)
	})

	t.Run("list second page", func(t *testing.T) {
		result, total, err := repo.List(ctx, 2, 2)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, result, 1)
	})
}

func TestUserRepository_EmailExists(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := createTestUser()
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	t.Run("existing email", func(t *testing.T) {
		exists, err := repo.EmailExists(ctx, user.Email)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("non-existing email", func(t *testing.T) {
		exists, err := repo.EmailExists(ctx, "nonexistent@example.com")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("empty email", func(t *testing.T) {
		_, err := repo.EmailExists(ctx, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})
}

func TestUserRepository_CountUsers(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create users
	users := []*models.User{
		{Email: "user1@example.com", Password: "password", FirstName: "User", LastName: "One", Role: models.RoleCustomer, IsActive: true},
		{Email: "user2@example.com", Password: "password", FirstName: "User", LastName: "Two", Role: models.RoleAdmin, IsActive: true},
		{Email: "user3@example.com", Password: "password", FirstName: "User", LastName: "Three", Role: models.RoleCustomer, IsActive: false},
	}

	for _, user := range users {
		err := repo.Create(ctx, user)
		require.NoError(t, err)
	}

	t.Run("count all users", func(t *testing.T) {
		count, err := repo.CountUsers(ctx)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})

	t.Run("count by role", func(t *testing.T) {
		count, err := repo.CountUsersByRole(ctx, models.RoleCustomer)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)

		count, err = repo.CountUsersByRole(ctx, models.RoleAdmin)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	t.Run("count active users", func(t *testing.T) {
		count, err := repo.CountActiveUsers(ctx)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)
	})
}

func TestUserRepository_ActivateDeactivate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := createTestUser()
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	t.Run("deactivate user", func(t *testing.T) {
		err := repo.DeactivateUser(ctx, user.ID)
		assert.NoError(t, err)

		updated, err := repo.GetByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.False(t, updated.IsActive)
	})

	t.Run("activate user", func(t *testing.T) {
		err := repo.ActivateUser(ctx, user.ID)
		assert.NoError(t, err)

		updated, err := repo.GetByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.True(t, updated.IsActive)
	})
}
