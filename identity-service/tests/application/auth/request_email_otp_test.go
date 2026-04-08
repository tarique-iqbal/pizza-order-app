package auth_test

import (
	"context"
	"errors"
	"identity-service/internal/application/auth"
	dAuth "identity-service/internal/domain/auth"
	"identity-service/internal/infrastructure/persistence"
	"identity-service/internal/infrastructure/security"
	"identity-service/internal/shared/event"
	"testing"

	"github.com/stretchr/testify/assert"
)

var emailOTP *auth.RequestEmailOTP
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

func requestEmailOTP() *auth.RequestEmailOTP {
	ts := testStorage()
	truncateTables(ts.DB)

	repo = persistence.NewEmailVerificationRepository(ts.DB)
	otp := security.NewOTPGenerator()
	mockPublisher = &MockEventPublisher{}

	return auth.NewRequestEmailOTP(repo, otp, mockPublisher)
}

func TestCreateEmailVerificationUseCase_Success(t *testing.T) {
	emailOTP = requestEmailOTP()

	input := auth.EmailVerificationRequest{
		Email: "adam.dangelo@example.com",
	}

	err := emailOTP.Execute(context.Background(), input)
	emailVerification, _ := repo.FindByEmail(context.Background(), input.Email)
	diff := emailVerification.ExpiresAt.Sub(emailVerification.CreatedAt)

	assert.Nil(t, err)
	assert.NotNil(t, emailVerification)
	assert.Equal(t, "adam.dangelo@example.com", emailVerification.Email)
	assert.InDelta(t, 15, diff.Minutes(), 0.001, "Delta threshold exceeded")

	createdEvent, ok := mockPublisher.PublishedEvents[0].(auth.EmailVerificationCreated)
	assert.True(t, ok)
	assert.Equal(t, "adam.dangelo@example.com", createdEvent.Email)
	assert.Equal(t, emailVerification.Code, createdEvent.Code)
	assert.Equal(t, "email.verification_created", createdEvent.GetEventName())
	assert.Len(t, mockPublisher.PublishedEvents, 1)
}
