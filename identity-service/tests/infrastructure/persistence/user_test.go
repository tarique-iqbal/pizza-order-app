package persistence_test

import (
	"context"
	"identity-service/internal/domain/user"
	"identity-service/internal/infrastructure/persistence"
	"identity-service/tests/infrastructure/db/fixtures"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupUserRepo() user.UserRepository {
	ts := testStorage()
	truncateTables(ts.DB)

	if err := fixtures.LoadUserFixtures(ts.DB); err != nil {
		panic(err)
	}

	return persistence.NewUserRepository(ts.DB)
}

func TestUserRepository_Create(t *testing.T) {
	userRepo := setupUserRepo()

	usr := user.User{
		FirstName: "Adam",
		LastName:  "D'Angelo",
		Email:     "adam.dangelo@example.com",
		Password:  "hashedpassword",
		Role:      "customer",
		CreatedAt: time.Now().UTC(),
	}

	userID, _ := uuid.NewV7()
	usr.ID = userID

	err := userRepo.Create(context.Background(), &usr)

	assert.Nil(t, err)
	assert.NotZero(t, usr.ID)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	userRepo := setupUserRepo()

	usr, err := userRepo.FindByEmail(context.Background(), "john.doe@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "John", usr.FirstName)
}

func TestUserRepository_EmailExists(t *testing.T) {
	userRepo := setupUserRepo()

	exists, err := userRepo.EmailExists("john.doe@example.com")
	assert.NoError(t, err)
	assert.True(t, exists, "Email is expected to be exists")

	exists, err = userRepo.EmailExists("random@example.com")
	assert.NoError(t, err)
	assert.False(t, exists, "Email is not expected to be exists")
}

func TestUserRepository_FindByID_Success(t *testing.T) {
	ctx := context.Background()
	userRepo := setupUserRepo()

	u := &user.User{
		FirstName: "Tony",
		LastName:  "Soprano",
		Email:     "tony@satrialis.com",
		Role:      "owner",
		Status:    "active",
		CreatedAt: time.Now().UTC(),
	}

	userID, _ := uuid.NewV7()
	u.ID = userID

	err := userRepo.Create(ctx, u)
	require.NoError(t, err)

	res, err := userRepo.FindByID(ctx, userID)

	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, u.ID, res.ID)
	assert.Equal(t, u.Email, res.Email)
	assert.Equal(t, u.FirstName, res.FirstName)
	assert.Equal(t, u.LastName, res.LastName)
	assert.Equal(t, u.Role, res.Role)
	assert.Equal(t, u.Status, res.Status)
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := setupUserRepo()

	userID, _ := uuid.NewV7()

	res, err := repo.FindByID(ctx, userID)

	require.NoError(t, err)
	assert.Nil(t, res)
}
