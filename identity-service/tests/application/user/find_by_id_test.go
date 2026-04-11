package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	app "identity-service/internal/application/user"
	domain "identity-service/internal/domain/user"
	"identity-service/internal/infrastructure/persistence"
	"identity-service/tests/infrastructure/db/fixtures"
)

func setupFindByID() *app.FindByID {
	ts := testStorage()
	truncateTables(ts.DB)

	userRepo := persistence.NewUserRepository(ts.DB)
	mockPublisher = &MockEventPublisher{}

	if err := fixtures.LoadUserFixtures(ts.DB); err != nil {
		panic(err)
	}

	return app.NewFindByID(userRepo)
}

func TestFindByID_Success(t *testing.T) {
	ctx := context.Background()
	ts := testStorage()
	uc := setupFindByID()

	u := &domain.User{
		FirstName: "Tony",
		LastName:  "Soprano",
		Email:     "tony@satrialis.com",
		Role:      "owner",
		Status:    "active",
		CreatedAt: time.Now(),
	}

	err := ts.DB.WithContext(ctx).Create(u).Error
	require.NoError(t, err)

	res, err := uc.Execute(ctx, 2)

	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, 2, res.ID)
	assert.Equal(t, "Tony", res.Name.First)
	assert.Equal(t, "Soprano", res.Name.Last)
	assert.Equal(t, "tony@satrialis.com", res.Email)
	assert.Equal(t, "owner", res.Role)
	assert.Equal(t, "active", res.Status)
}

func TestFindByID_NotFound(t *testing.T) {
	ctx := context.Background()
	uc := setupFindByID()

	res, err := uc.Execute(ctx, 767)

	require.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Equal(t, app.Response{}, res)
}
