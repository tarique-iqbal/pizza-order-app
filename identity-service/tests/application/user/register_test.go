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
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockPublisher *MockEventPublisher

type MockEventPublisher struct {
	PublishedEvents []event.Event
	ShouldFail      bool
}

func (m *MockEventPublisher) Publish(e event.Event) error {
	if m.ShouldFail {
		return errors.New("mock publish failure")
	}
	m.PublishedEvents = append(m.PublishedEvents, e)
	return nil
}

func setupRegister() *user.Register {
	ts := testStorage()
	truncateTables(ts.DB)

	emailVerificationRepo := persistence.NewEmailVerificationRepository(ts.DB)
	codeVerifier := auth.NewEmailVerifier(emailVerificationRepo)
	userRepo := persistence.NewUserRepository(ts.DB)
	hasher := security.NewPasswordHasher()
	mockPublisher = &MockEventPublisher{}

	if err := fixtures.LoadEmailVerificationFixtures(ts.DB); err != nil {
		panic(err)
	}
	if err := fixtures.LoadUserFixtures(ts.DB); err != nil {
		panic(err)
	}

	return user.NewRegister(codeVerifier, userRepo, hasher, mockPublisher)
}

func TestRegister(t *testing.T) {
	register := setupRegister()

	input := user.RegisterRequest{
		FirstName: "Adam",
		LastName:  "D'Angelo",
		Email:     "adam.dangelo@example.com",
		Password:  "securepassword",
		Code:      "476190",
	}

	newUser, err := register.Execute(context.Background(), input)

	assert.Nil(t, err)
	assert.NotNil(t, newUser)
	assert.Equal(t, "Adam", newUser.FirstName)

	createdEvent, ok := mockPublisher.PublishedEvents[0].(user.UserRegistered)
	assert.True(t, ok)
	assert.Equal(t, "Adam", createdEvent.FirstName)
	assert.Equal(t, "adam.dangelo@example.com", createdEvent.Email)
	assert.Equal(t, "user.registered", createdEvent.GetEventName())
	assert.Len(t, mockPublisher.PublishedEvents, 1)
}
