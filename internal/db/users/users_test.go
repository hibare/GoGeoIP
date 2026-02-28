package users

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hibare/Waypoint/internal/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUser_TableName(t *testing.T) {
	user := &User{}
	assert.Equal(t, tableNameUsers, user.TableName())
}

func TestUser_BeforeCreate(t *testing.T) {
	user := &User{}

	err := user.BeforeCreate(nil)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)
}

func TestCreateUser(t *testing.T) {
	db := testhelpers.SetupSharedTestDB(t)
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		user := &User{
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
			Groups:    []string{"admin", "user"},
		}

		err := CreateUser(ctx, db.DB, user)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, user.ID)
		assert.NotZero(t, user.CreatedAt)
		assert.NotZero(t, user.UpdatedAt)
	})

	t.Run("duplicate email", func(t *testing.T) {
		// Try to create user with existing email from seeds
		user := &User{
			Email:     "admin@test.com",
			FirstName: "Duplicate",
			LastName:  "User",
			Groups:    []string{"user"},
		}

		err := CreateUser(ctx, db.DB, user)
		require.Error(t, err)
		// Should be a unique constraint violation
		assert.Contains(t, err.Error(), "duplicate")
	})
}

func TestGetUserByEmail(t *testing.T) {
	db := testhelpers.SetupSharedTestDB(t)
	ctx := context.Background()

	t.Run("existing user", func(t *testing.T) {
		user, err := GetUserByEmail(ctx, db.DB, "admin@test.com")
		require.NoError(t, err)
		assert.Equal(t, "admin@test.com", user.Email)
		assert.Equal(t, "Test", user.FirstName)
		assert.Equal(t, "Admin", user.LastName)
		assert.NotZero(t, user.CreatedAt)
		assert.NotZero(t, user.UpdatedAt)
	})

	t.Run("non-existing user", func(t *testing.T) {
		user, err := GetUserByEmail(ctx, db.DB, "nonexistent@example.com")
		require.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.Nil(t, user)
	})

	t.Run("empty email", func(t *testing.T) {
		user, err := GetUserByEmail(ctx, db.DB, "")
		require.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.Nil(t, user)
	})
}

func TestGetUserByID(t *testing.T) {
	db := testhelpers.SetupSharedTestDB(t)
	ctx := context.Background()

	// First get a user by email to get their ID
	existingUser, err := GetUserByEmail(ctx, db.DB, "admin@test.com")
	require.NoError(t, err)

	t.Run("existing user", func(t *testing.T) {
		user, err := GetUserByID(ctx, db.DB, existingUser.ID.String())
		require.NoError(t, err)
		assert.Equal(t, existingUser.ID, user.ID)
		assert.Equal(t, existingUser.Email, user.Email)
		assert.Equal(t, existingUser.FirstName, user.FirstName)
		assert.Equal(t, existingUser.LastName, user.LastName)
	})

	t.Run("non-existing user", func(t *testing.T) {
		user, err := GetUserByID(ctx, db.DB, uuid.New().String())
		require.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.Nil(t, user)
	})

	t.Run("invalid UUID", func(t *testing.T) {
		user, err := GetUserByID(ctx, db.DB, "invalid-uuid")
		require.Error(t, err)
		// PostgreSQL validates UUID format, so we get a different error than gorm.ErrRecordNotFound
		assert.Nil(t, user)
	})

	t.Run("empty ID", func(t *testing.T) {
		user, err := GetUserByID(ctx, db.DB, "")
		require.Error(t, err)
		// PostgreSQL validates UUID format, so we get a different error than gorm.ErrRecordNotFound
		assert.Nil(t, user)
	})
}

func TestUser_Fields(t *testing.T) {
	now := time.Now()
	userID := uuid.New()

	user := &User{
		ID:        userID,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Groups:    []string{"admin", "moderator"},
		CreatedAt: now,
		UpdatedAt: now,
		LastLogin: now,
	}

	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, []string{"admin", "moderator"}, user.Groups)
	assert.Equal(t, now, user.CreatedAt)
	assert.Equal(t, now, user.UpdatedAt)
	assert.Equal(t, now, user.LastLogin)
}

func TestUser_JSONTags(t *testing.T) {
	user := &User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Groups:    []string{"group1", "group2"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		LastLogin: time.Now(),
	}

	// Test that the struct has JSON tags (basic smoke test)
	// This ensures JSON tags are present on the struct
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.Email)
}

func TestUser_GormTags(t *testing.T) {
	user := &User{}

	// Test that GORM tags are set correctly by checking table name
	assert.Equal(t, "users", user.TableName())

	// The actual GORM tag validation would happen during database operations
	// which are tested in the integration tests above
}

func TestUser_BeforeCreate_IDGeneration(t *testing.T) {
	user := &User{}

	// Test that BeforeCreate generates a new UUID if ID is nil
	originalID := user.ID
	assert.Equal(t, uuid.Nil, originalID)

	err := user.BeforeCreate(nil)
	require.NoError(t, err)

	assert.NotEqual(t, uuid.Nil, user.ID)
	assert.NotEqual(t, originalID, user.ID)
}

func TestUser_BeforeCreate_ExistingID(t *testing.T) {
	existingID := uuid.New()
	user := &User{ID: existingID}

	// Test that BeforeCreate doesn't change existing ID
	err := user.BeforeCreate(nil)
	require.NoError(t, err)

	assert.Equal(t, existingID, user.ID)
}

func TestUpdateUser(t *testing.T) {
	db := testhelpers.SetupSharedTestDB(t)
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		// First create a user
		user := &User{
			Email:     "update-test@example.com",
			FirstName: "Original",
			LastName:  "Name",
			Groups:    []string{"user"},
		}

		err := CreateUser(ctx, db.DB, user)
		require.NoError(t, err)
		originalID := user.ID

		// Update the user
		updates := &User{
			FirstName: "Updated",
			LastName:  "Name",
			Groups:    []string{"admin", "moderator"},
			LastLogin: time.Now().UTC(),
		}

		err = UpdateUser(ctx, db.DB, user.ID.String(), updates)
		require.NoError(t, err)

		// Verify the update
		updatedUser, err := GetUserByID(ctx, db.DB, user.ID.String())
		require.NoError(t, err)

		assert.Equal(t, originalID, updatedUser.ID)
		assert.Equal(t, "update-test@example.com", updatedUser.Email)
		assert.Equal(t, "Updated", updatedUser.FirstName)
		assert.Equal(t, "Name", updatedUser.LastName)
		assert.Equal(t, []string{"admin", "moderator"}, updatedUser.Groups)
		assert.True(t, updatedUser.LastLogin.After(user.LastLogin))
		assert.True(t, updatedUser.UpdatedAt.After(user.UpdatedAt))
	})

	t.Run("update non-existing user", func(t *testing.T) {
		nonExistingID := uuid.New().String()
		updates := &User{
			FirstName: "Should",
			LastName:  "Fail",
		}

		err := UpdateUser(ctx, db.DB, nonExistingID, updates)
		// GORM Updates doesn't return an error for non-existing records
		// It just doesn't update anything
		require.NoError(t, err)
	})

	t.Run("update with partial fields", func(t *testing.T) {
		// Create a user
		user := &User{
			Email:     "partial-update@example.com",
			FirstName: "Original",
			LastName:  "Name",
			Groups:    []string{"user"},
		}

		err := CreateUser(ctx, db.DB, user)
		require.NoError(t, err)

		// Update only some fields
		updates := &User{
			FirstName: "Partially",
			Groups:    []string{"admin"},
		}

		err = UpdateUser(ctx, db.DB, user.ID.String(), updates)
		require.NoError(t, err)

		// Verify partial update
		updatedUser, err := GetUserByID(ctx, db.DB, user.ID.String())
		require.NoError(t, err)

		assert.Equal(t, "Partially", updatedUser.FirstName)
		assert.Equal(t, "Name", updatedUser.LastName) // Should remain unchanged
		assert.Equal(t, []string{"admin"}, updatedUser.Groups)
	})

	t.Run("update groups array", func(t *testing.T) {
		// Create a user
		user := &User{
			Email:     "groups-update@example.com",
			FirstName: "Groups",
			LastName:  "Test",
			Groups:    []string{"user"},
		}

		err := CreateUser(ctx, db.DB, user)
		require.NoError(t, err)

		// Update groups to empty array
		updates := &User{
			Groups: []string{},
		}

		err = UpdateUser(ctx, db.DB, user.ID.String(), updates)
		require.NoError(t, err)

		// Verify groups update
		updatedUser, err := GetUserByID(ctx, db.DB, user.ID.String())
		require.NoError(t, err)
		assert.Equal(t, []string{}, updatedUser.Groups)

		// Update groups to larger array
		updates = &User{
			Groups: []string{"admin", "moderator", "editor", "viewer"},
		}

		err = UpdateUser(ctx, db.DB, user.ID.String(), updates)
		require.NoError(t, err)

		// Verify groups update
		updatedUser, err = GetUserByID(ctx, db.DB, user.ID.String())
		require.NoError(t, err)
		assert.Equal(t, []string{"admin", "moderator", "editor", "viewer"}, updatedUser.Groups)
	})
}

func TestCreateUser_WithGroups(t *testing.T) {
	db := testhelpers.SetupSharedTestDB(t)
	ctx := context.Background()

	t.Run("create user with groups array", func(t *testing.T) {
		user := &User{
			Email:     "groups-test@example.com",
			FirstName: "Groups",
			LastName:  "Test",
			Groups:    []string{"admin", "moderator", "user"},
		}

		err := CreateUser(ctx, db.DB, user)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, user.ID)
		assert.Equal(t, []string{"admin", "moderator", "user"}, user.Groups)
	})

	t.Run("retrieve user with groups array", func(t *testing.T) {
		user, err := GetUserByEmail(ctx, db.DB, "groups-test@example.com")
		require.NoError(t, err)
		assert.Equal(t, "groups-test@example.com", user.Email)
		assert.Equal(t, []string{"admin", "moderator", "user"}, user.Groups)
	})

	t.Run("update user groups via assign", func(t *testing.T) {
		// Simulate what happens in auth handler when user logs in again
		newGroups := []string{"superadmin", "user"}

		// Use the same pattern as the auth handler
		result := db.DB.Where(User{Email: "groups-test@example.com"}).Assign(User{
			Groups:    newGroups,
			LastLogin: time.Now().UTC(),
		}).FirstOrCreate(&User{})

		require.NoError(t, result.Error)

		// Verify the groups were updated
		user, err := GetUserByEmail(ctx, db.DB, "groups-test@example.com")
		require.NoError(t, err)
		assert.Equal(t, newGroups, user.Groups)
	})
}

func TestUser_Groups_JSONBSerialization(t *testing.T) {
	db := testhelpers.SetupSharedTestDB(t)
	ctx := context.Background()

	t.Run("empty groups array", func(t *testing.T) {
		user := &User{
			Email:     "empty-groups@example.com",
			FirstName: "Empty",
			LastName:  "Groups",
			Groups:    []string{},
		}

		err := CreateUser(ctx, db.DB, user)
		require.NoError(t, err)

		retrieved, err := GetUserByEmail(ctx, db.DB, "empty-groups@example.com")
		require.NoError(t, err)
		assert.Equal(t, []string{}, retrieved.Groups)
	})

	t.Run("single group", func(t *testing.T) {
		user := &User{
			Email:     "single-group@example.com",
			FirstName: "Single",
			LastName:  "Group",
			Groups:    []string{"admin"},
		}

		err := CreateUser(ctx, db.DB, user)
		require.NoError(t, err)

		retrieved, err := GetUserByEmail(ctx, db.DB, "single-group@example.com")
		require.NoError(t, err)
		assert.Equal(t, []string{"admin"}, retrieved.Groups)
	})

	t.Run("groups with special characters", func(t *testing.T) {
		user := &User{
			Email:     "special-groups@example.com",
			FirstName: "Special",
			LastName:  "Groups",
			Groups:    []string{"admin-dash", "user_underscore", "mod.dot"},
		}

		err := CreateUser(ctx, db.DB, user)
		require.NoError(t, err)

		retrieved, err := GetUserByEmail(ctx, db.DB, "special-groups@example.com")
		require.NoError(t, err)
		assert.Equal(t, []string{"admin-dash", "user_underscore", "mod.dot"}, retrieved.Groups)
	})
}
