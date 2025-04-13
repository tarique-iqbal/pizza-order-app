package user_test

import (
	"errors"
	aUser "pizza-order-api/internal/application/user"
	"pizza-order-api/internal/infrastructure/persistence"
	"pizza-order-api/internal/shared/event"
	"pizza-order-api/tests/internal/infrastructure/db"
	"pizza-order-api/tests/internal/infrastructure/db/fixtures"
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
	userRepo := persistence.NewUserRepository(testDB)
	mockPublisher = &MockEventPublisher{}

	if err := fixtures.LoadUserFixtures(testDB); err != nil {
		panic(err)
	}

	return aUser.NewCreateUserUseCase(userRepo, mockPublisher)
}

func TestCreateUserUseCase(t *testing.T) {
	createUseCase = createUserUseCase()

	input := aUser.UserCreateDTO{
		FirstName: "Adam",
		LastName:  "D'Angelo",
		Email:     "adam.dangelo@example.com",
		Password:  "securepassword",
	}

	user, err := createUseCase.Execute(input)

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Adam", user.FirstName)

	createdEvent, ok := mockPublisher.PublishedEvents[0].(aUser.UserCreatedEvent)
	assert.True(t, ok)
	assert.Equal(t, "Adam", createdEvent.FirstName)
	assert.Equal(t, "D'Angelo", createdEvent.LastName)
	assert.Equal(t, "adam.dangelo@example.com", createdEvent.Email)
	assert.Equal(t, "user.registered", createdEvent.GetEventName())
	assert.Len(t, mockPublisher.PublishedEvents, 1)
}
