package user_test

import (
	"context"
	"identity-service/internal/application/user"
	"identity-service/internal/infrastructure/auth"
	"identity-service/internal/infrastructure/persistence"
	"identity-service/internal/infrastructure/security"
	"identity-service/tests/infrastructure/db/fixtures"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupRegisterOwner() *user.RegisterOwner {
	ts := testStorage()
	truncateTables(ts.DB)

	emailVerificationRepo := persistence.NewEmailVerificationRepository(ts.DB)
	codeVerifier := auth.NewEmailVerifier(emailVerificationRepo)
	userRepo := persistence.NewUserRepository(ts.DB)
	outboxRepo := persistence.NewOutboxRepository(ts.DB)
	hasher := security.NewPasswordHasher()
	mockPublisher = &MockEventPublisher{}

	if err := fixtures.LoadEmailVerificationFixtures(ts.DB); err != nil {
		panic(err)
	}
	if err := fixtures.LoadUserFixtures(ts.DB); err != nil {
		panic(err)
	}

	return user.NewRegisterOwner(ts.DB, codeVerifier, hasher, userRepo, outboxRepo, mockPublisher)
}

func TestRegisterOwner_Success(t *testing.T) {
	registerOwner := setupRegisterOwner()

	input := user.RegisterOwnerRequest{
		RegisterCustomerRequest: user.RegisterCustomerRequest{
			FirstName: "Sophie",
			LastName:  "Müller",
			Email:     "sophie.mueller@example.com",
			Password:  "securepassword",
			Role:      "owner",
			Code:      "365189",
		},
		BusinessName: "Domino's Pizza",
		VATNumber:    "DE987654321",
	}

	newUser, err := registerOwner.Execute(context.Background(), input)

	assert.Nil(t, err)
	assert.NotNil(t, newUser)
	assert.Equal(t, "Sophie", newUser.Name.First)

	createdEvent, ok := mockPublisher.PublishedEvents[0].(user.UserRegistered)
	assert.True(t, ok)
	assert.Equal(t, "Sophie", createdEvent.FirstName)
	assert.Equal(t, "sophie.mueller@example.com", createdEvent.Email)
	assert.Equal(t, "user.registered", createdEvent.GetEventName())
	assert.Len(t, mockPublisher.PublishedEvents, 1)
}

func TestRegisterOwner_Failure_EmailVerification(t *testing.T) {
	registerOwner := setupRegisterOwner()

	input := user.RegisterOwnerRequest{
		RegisterCustomerRequest: user.RegisterCustomerRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "invalid@example.com",
			Password:  "password",
			Role:      "owner",
			Code:      "wrong-code", // invalid
		},
		BusinessName: "Test Biz",
		VATNumber:    "DE111654321",
	}

	res, err := registerOwner.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Empty(t, res)

	assert.Len(t, mockPublisher.PublishedEvents, 0)
}

func TestRegisterOwner_Failure_DuplicateEmail(t *testing.T) {
	registerOwner := setupRegisterOwner()

	input := user.RegisterOwnerRequest{
		RegisterCustomerRequest: user.RegisterCustomerRequest{
			FirstName: "Existing",
			LastName:  "User",
			Email:     "existing@example.com", // from fixture
			Password:  "password",
			Role:      "owner",
			Code:      "365189",
		},
		BusinessName: "Biz",
		VATNumber:    "DE222654321",
	}

	res, err := registerOwner.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Empty(t, res)

	assert.Len(t, mockPublisher.PublishedEvents, 0)
}

func TestRegisterOwner_PublishFails_ShouldStillSucceed(t *testing.T) {
	registerOwner := setupRegisterOwner()

	// override publisher to fail
	mockPublisher.ShouldFail = true

	input := user.RegisterOwnerRequest{
		RegisterCustomerRequest: user.RegisterCustomerRequest{
			FirstName: "Alice",
			LastName:  "Schmidt",
			Email:     "alice@example.com",
			Password:  "password",
			Role:      "owner",
			Code:      "347578",
		},
		BusinessName: "Pizza Hub",
		VATNumber:    "DE444654321",
	}

	res, err := registerOwner.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotEmpty(t, res)

	// event attempted
	assert.Len(t, mockPublisher.PublishedEvents, 1)
}
