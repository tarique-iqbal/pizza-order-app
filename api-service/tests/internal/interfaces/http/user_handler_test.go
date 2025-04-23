package http_test

import (
	"api-service/internal/application/user"
	"api-service/internal/infrastructure/persistence"
	"api-service/internal/infrastructure/security"
	iValidator "api-service/internal/infrastructure/validator"
	uiHttp "api-service/internal/interfaces/http"
	"api-service/internal/shared/event"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
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

func setupUserHandler() *uiHttp.UserHandler {
	userRepo := persistence.NewUserRepository(testDB)
	hasher := security.NewPasswordHasher()
	mockPublisher = &MockEventPublisher{}

	createUserUseCase := user.NewCreateUserUseCase(userRepo, hasher, mockPublisher)
	customValidator := iValidator.NewCustomValidator(userRepo, nil)
	userUseCases := &uiHttp.UserUseCases{
		CreateUser:      createUserUseCase,
		CustomValidator: customValidator,
	}

	return uiHttp.NewUserHandler(userUseCases)
}

func TestUserHandler_CreateUser_Success(t *testing.T) {
	uHandler := setupUserHandler()
	router := gin.Default()
	router.POST("/api/users", uHandler.Create)

	reqBody, _ := json.Marshal(map[string]string{
		"first_name": "Alice",
		"last_name":  "Doe",
		"email":      "alice@example.com",
		"password":   "pass123",
		"role":       "User",
	})

	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "alice@example.com")
}
