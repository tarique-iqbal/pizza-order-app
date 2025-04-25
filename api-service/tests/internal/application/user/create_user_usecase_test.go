package user_test

import (
	aUser "api-service/internal/application/user"
	"api-service/internal/infrastructure/auth"
	"api-service/internal/infrastructure/persistence"
	"api-service/internal/infrastructure/security"
	"api-service/internal/shared/event"
	"api-service/tests/internal/infrastructure/db"
	"api-service/tests/internal/infrastructure/db/fixtures"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var createUseCase *aUser.CreateUserUseCase
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
	createUseCase = createUserUseCase()

	input := aUser.UserCreateDTO{
		FirstName: "Adam",
		LastName:  "D'Angelo",
		Email:     "adam.dangelo@example.com",
		Password:  "securepassword",
		Code:      "476190",
	}

	user, err := createUseCase.Execute(input)

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
