package http_test

import (
	"api-service/internal/application/auth"
	"api-service/internal/domain/user"
	"api-service/internal/infrastructure/persistence"
	"api-service/internal/infrastructure/security"
	uiHttp "api-service/internal/interfaces/http"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAuthHandler() *uiHttp.AuthHandler {
	userRepo := persistence.NewUserRepository(testDB)
	hasher := security.NewPasswordHasher()
	jwt := security.NewJWTService("TestSecretKey")

	signInUseCase := auth.NewSignInUseCase(userRepo, hasher, jwt)
	authUseCases := &uiHttp.AuthUseCases{
		SignIn: signInUseCase,
	}

	return uiHttp.NewAuthHandler(authUseCases)
}

func TestAuthHandler_SignIn_Success(t *testing.T) {
	aHandler := setupAuthHandler()
	repo := persistence.NewUserRepository(testDB)
	hasher := security.NewPasswordHasher()
	jwt := security.NewJWTService("TestSecretKey")
	hp, _ := hasher.Hash("password123")

	newUser := &user.User{
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
		router.POST("/api/auth/signin", aHandler.SignIn)

		body, _ := json.Marshal(tc.requestBody)
		req, _ := http.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		jsonString := recorder.Body.String()
		var data map[string]interface{}
		json.Unmarshal([]byte(jsonString), &data)
		_, err := jwt.ParseToken(data["token"].(string))

		assert.NoError(t, err)
		assert.Equal(t, tc.expectedCode, recorder.Code)
	})
}

func TestAuthHandler_SignIn_Failed(t *testing.T) {
	aHandler := setupAuthHandler()
	repo := persistence.NewUserRepository(testDB)
	hasher := security.NewPasswordHasher()
	hp, _ := hasher.Hash("password123")

	newUser := &user.User{
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
			router.POST("/api/auth/signin", aHandler.SignIn)

			body, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedCode, recorder.Code)
			assert.JSONEq(t, tc.expectedBody, recorder.Body.String())
		})
	}
}
