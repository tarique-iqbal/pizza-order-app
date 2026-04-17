package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	userapp "identity-service/internal/application/user"
	"identity-service/internal/domain/user"
	"identity-service/internal/infrastructure/persistence"
	"identity-service/tests/infrastructure/db/fixtures"
)

func setupFindByID() *userapp.FindByID {
	ts := testStorage()
	truncateTables(ts.DB)

	userRepo := persistence.NewUserRepository(ts.DB)
	mockPublisher = &MockEventPublisher{}

	if err := fixtures.LoadUserFixtures(ts.DB); err != nil {
		panic(err)
	}

	return userapp.NewFindByID(userRepo)
}

func TestFindByID_Success(t *testing.T) {
	ctx := context.Background()
	ts := testStorage()
	uc := setupFindByID()

	u := &user.User{
		FirstName: "Tony",
		LastName:  "Soprano",
		Email:     "tony@satrialis.com",
		Role:      "owner",
		Status:    "active",
		CreatedAt: time.Now(),
	}

	userID, _ := uuid.NewV7()
	u.ID = userID

	err := ts.DB.WithContext(ctx).Create(u).Error
	require.NoError(t, err)

	res, err := uc.Execute(ctx, userID)

	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, userID, res.ID)
	assert.Equal(t, "Tony", res.Name.First)
	assert.Equal(t, "Soprano", res.Name.Last)
	assert.Equal(t, "tony@satrialis.com", res.Email)
	assert.Equal(t, "owner", res.Role)
	assert.Equal(t, "active", res.Status)
}

func TestFindByID_NotFound(t *testing.T) {
	ctx := context.Background()
	uc := setupFindByID()

	userID, _ := uuid.NewV7()

	res, err := uc.Execute(ctx, userID)

	require.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Equal(t, userapp.Response{}, res)
}
