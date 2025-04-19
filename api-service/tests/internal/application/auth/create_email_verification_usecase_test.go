package auth_test

import (
	"api-service/internal/application/auth"
	dAuth "api-service/internal/domain/auth"
	"api-service/internal/infrastructure/persistence"
	"api-service/internal/infrastructure/security"
	"api-service/internal/shared/event"
	"api-service/tests/internal/infrastructure/db"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var createUseCase *auth.CreateEmailVerificationUseCase
var mockPublisher *MockEventPublisher
var repo dAuth.EmailVerificationRepository

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

func createEmailVerificationUseCase() *auth.CreateEmailVerificationUseCase {
	testDB := db.SetupTestDB()
	repo = persistence.NewEmailVerificationRepository(testDB)
	otp := security.NewSixDigitOTPGenerator()
	mockPublisher = &MockEventPublisher{}

	return auth.NewCreateEmailVerificationUseCase(repo, otp, mockPublisher)
}

func TestCreateEmailVerificationUseCase_Success(t *testing.T) {
	createUseCase = createEmailVerificationUseCase()

	input := auth.EmailVerificationRequestDTO{
		Email: "adam.dangelo@example.com",
	}

	err := createUseCase.Execute(input)
	emailVerification, _ := repo.FindByEmail(input.Email)

	assert.Nil(t, err)
	assert.NotNil(t, emailVerification)
	assert.Equal(t, "adam.dangelo@example.com", emailVerification.Email)

	createdEvent, ok := mockPublisher.PublishedEvents[0].(auth.EmailVerificationCreatedEvent)
	assert.True(t, ok)
	assert.Equal(t, "adam.dangelo@example.com", createdEvent.Email)
	assert.Equal(t, emailVerification.Code, createdEvent.Code)
	assert.Equal(t, "email.verification_created", createdEvent.GetEventName())
	assert.Len(t, mockPublisher.PublishedEvents, 1)
}
