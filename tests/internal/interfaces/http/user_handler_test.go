package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	aUser "pizza-order-api/internal/application/user"
	dUser "pizza-order-api/internal/domain/user"
	"pizza-order-api/internal/infrastructure/persistence"
	"pizza-order-api/internal/infrastructure/security"
	iValidator "pizza-order-api/internal/infrastructure/validator"
	uiHttp "pizza-order-api/internal/interfaces/http"
	"pizza-order-api/internal/shared/event"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var uHandler *uiHttp.UserHandler
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
	mockPublisher = &MockEventPublisher{}

	createUserUseCase := aUser.NewCreateUserUseCase(userRepo, mockPublisher)
	signInUserUseCase := aUser.NewSignInUserUseCase(userRepo)
	customValidator := iValidator.NewCustomValidator(userRepo, nil)
	userUseCases := &uiHttp.UserUseCases{
		CreateUser:      createUserUseCase,
		SignIn:          signInUserUseCase,
		CustomValidator: customValidator,
	}

	return uiHttp.NewUserHandler(userUseCases)
}

func TestUserHandler_CreateUser_Success(t *testing.T) {
	uHandler = setupUserHandler()
	router := gin.Default()
	router.POST("/api/users/signup", uHandler.CreateUser)

	reqBody, _ := json.Marshal(map[string]string{
		"first_name": "Alice",
		"last_name":  "Doe",
		"email":      "alice@example.com",
		"password":   "pass123",
		"role":       "user",
	})

	req, _ := http.NewRequest("POST", "/api/users/signup", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "alice@example.com")
}

func TestUserHandler_SignIn_Success(t *testing.T) {
	uHandler = setupUserHandler()
	repo := persistence.NewUserRepository(testDB)
	hp, _ := security.HashPassword("password123")

	newUser := &dUser.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@example.com",
		Password:  hp,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	repo.Create(newUser)

	tc := struct {
		name         string
		requestBody  map[string]string
		expectedCode int
		expectedBody string
	}{
		name: "Successful Login",
		requestBody: map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		},
		expectedCode: http.StatusOK,
	}

	t.Run(tc.name, func(t *testing.T) {
		router := gin.Default()
		router.POST("/api/users/signin", uHandler.SignIn)

		body, _ := json.Marshal(tc.requestBody)
		req, _ := http.NewRequest(http.MethodPost, "/api/users/signin", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		jsonString := recorder.Body.String()
		var data map[string]interface{}
		json.Unmarshal([]byte(jsonString), &data)
		_, err := security.ValidateToken(data["token"].(string))

		assert.NoError(t, err)
		assert.Equal(t, tc.expectedCode, recorder.Code)
	})
}

func TestUserHandler_SignIn_Failed(t *testing.T) {
	uHandler = setupUserHandler()
	repo := persistence.NewUserRepository(testDB)
	hp, _ := security.HashPassword("password123")

	newUser := &dUser.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "test@example.com",
		Password:  hp,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	repo.Create(newUser)

	tests := []struct {
		name         string
		requestBody  map[string]string
		expectedCode int
		expectedBody string
	}{
		{
			name: "Successful Login",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"invalid credentials"}`,
		},
		{
			name: "Invalid Credentials",
			requestBody: map[string]string{
				"email":    "no.user.found@example.com",
				"password": "random",
			},
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"no record found"}`,
		},
		{
			name: "Invalid Request Body",
			requestBody: map[string]string{
				"email":    "invalid-email.com",
				"password": "password123",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"errors":[{"field":"Email","message":"Please provide a valid email address."}]}`,
		},
		{
			name: "Invalid Request Body",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"errors":[{"field":"Password","message":"This field is required."}]}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.Default()
			router.POST("/api/users/signin", uHandler.SignIn)

			body, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/users/signin", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedCode, recorder.Code)
			assert.JSONEq(t, tc.expectedBody, recorder.Body.String())
		})
	}
}
