package user_test

import (
	"context"
	"errors"
	aUser "identity-service/internal/application/user"
	"identity-service/internal/infrastructure/auth"
	"identity-service/internal/infrastructure/persistence"
	"identity-service/internal/infrastructure/security"
	"identity-service/internal/shared/event"
	"identity-service/tests/infrastructure/db"
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

func createUserUseCase() *aUser.CreateUserUseCase {
	testDB := db.SetupTestDB()

	emailVerificationRepo := persistence.NewEmailVerificationRepository(testDB)
	codeVerifier := auth.NewCodeVerificationService(emailVerificationRepo)
	userRepo := persistence.NewUserRepository(testDB)
	hasher := security.NewPasswordHasher()
	mockPublisher = &MockEventPublisher{}

	if err := fixtures.LoadEmailVerificationFixtures(testDB); err != nil {
		panic(err)
	}
	if err := fixtures.LoadUserFixtures(testDB); err != nil {
		panic(err)
	}

	return aUser.NewCreateUserUseCase(codeVerifier, userRepo, hasher, mockPublisher)
}

func TestCreateUserUseCase(t *testing.T) {
	createUserUC := createUserUseCase()

	input := aUser.UserCreateDTO{
		FirstName: "Adam",
		LastName:  "D'Angelo",
		Email:     "adam.dangelo@example.com",
		Password:  "securepassword",
		Code:      "476190",
	}

	user, err := createUserUC.Execute(context.Background(), input)

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Adam", user.FirstName)

	createdEvent, ok := mockPublisher.PublishedEvents[0].(aUser.UserCreatedEvent)
	assert.True(t, ok)
	assert.Equal(t, "Adam", createdEvent.FirstName)
	assert.Equal(t, "adam.dangelo@example.com", createdEvent.Email)
	assert.Equal(t, "user.registered", createdEvent.GetEventName())
	assert.Len(t, mockPublisher.PublishedEvents, 1)
}
