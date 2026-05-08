package user_test

import (
	"context"
	"errors"
	"identity-service/internal/application/user"
	"identity-service/internal/infrastructure/auth"
	"identity-service/internal/infrastructure/persistence"
	"identity-service/internal/infrastructure/security"
	"identity-service/internal/shared/event"
	"identity-service/tests/infrastructure/db/fixtures"
	"identity-service/tests/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockPublisher *MockEventPublisher

type MockEventPublisher struct {
	PublishedEvents []event.Event
	PublishedRaw    [][]byte
	ShouldFail      bool
}

func (m *MockEventPublisher) PublishEvent(ctx context.Context, e event.Event) error {
	m.PublishedEvents = append(m.PublishedEvents, e)
	if m.ShouldFail {
		return errors.New("mock publish failure")
	}
	return nil
}

func (m *MockEventPublisher) PublishRaw(ctx context.Context, topic string, jsonData []byte) error {
	m.PublishedRaw = append(m.PublishedRaw, jsonData)
	if m.ShouldFail {
		return errors.New("mock raw publish failure")
	}
	return nil
}

func setupRegisterCustomer(t *testing.T) *user.RegisterCustomer {
	db := testutil.DB(t)
	db.TruncateTables(t, testutil.TableEmailVerification, testutil.TableUser)

	_ = fixtures.LoadEmailVerificationFixtures(t, db.DB)
	_ = fixtures.LoadUserFixtures(t, db.DB)

	emailVerificationRepo := persistence.NewEmailVerificationRepository(db.DB)
	codeVerifier := auth.NewEmailVerifier(emailVerificationRepo)
	userRepo := persistence.NewUserRepository(db.DB)
	hasher := security.NewPasswordHasher()
	mockPublisher = &MockEventPublisher{}

	return user.NewRegisterCustomer(codeVerifier, userRepo, hasher, mockPublisher)
}

func TestRegisterCustomer_Success(t *testing.T) {
	register := setupRegisterCustomer(t)

	input := user.RegisterCustomerRequest{
		FirstName: "Adam",
		LastName:  "D'Angelo",
		Email:     "adam.dangelo@example.com",
		Password:  "securepassword",
		Code:      "476190", // from fixture
	}

	newUser, err := register.Execute(context.Background(), input)

	assert.Nil(t, err)
	assert.NotNil(t, newUser)
	assert.Equal(t, "Adam", newUser.Name.First)

	createdEvent, ok := mockPublisher.PublishedEvents[0].(user.UserRegistered)
	assert.True(t, ok)
	assert.Equal(t, "Adam", createdEvent.FirstName)
	assert.Equal(t, "adam.dangelo@example.com", createdEvent.Email)
	assert.Equal(t, "user.registered", createdEvent.GetEventName())
	assert.Len(t, mockPublisher.PublishedEvents, 1)
}

func TestRegisterCustomer_Failure_EmailVerification(t *testing.T) {
	register := setupRegisterCustomer(t)

	input := user.RegisterCustomerRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "invalid@example.com",
		Password:  "password",
		Code:      "wrong-code", // invalid
	}

	response, err := register.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Empty(t, response)

	assert.Len(t, mockPublisher.PublishedEvents, 0)
}

func TestRegisterCustomer_Failure_DuplicateEmail(t *testing.T) {
	register := setupRegisterCustomer(t)

	input := user.RegisterCustomerRequest{
		FirstName: "Existing",
		LastName:  "User",
		Email:     "existing@example.com", // from fixture
		Password:  "password",
		Code:      "365189",
	}

	response, err := register.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Empty(t, response)

	assert.Len(t, mockPublisher.PublishedEvents, 0)
}

func TestRegisterCustomer_PublishFails_ShouldStillSucceed(t *testing.T) {
	register := setupRegisterCustomer(t)

	// override publisher to fail
	mockPublisher.ShouldFail = true

	input := user.RegisterCustomerRequest{
		FirstName: "Alice",
		LastName:  "Schmidt",
		Email:     "alice@example.com",
		Password:  "password",
		Code:      "347578",
	}

	response, err := register.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotEmpty(t, response)

	// event attempted
	assert.Len(t, mockPublisher.PublishedEvents, 1)
}
