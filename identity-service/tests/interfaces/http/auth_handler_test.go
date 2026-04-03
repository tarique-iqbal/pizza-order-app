package http_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"identity-service/internal/application/auth"
	"identity-service/internal/domain/user"
	"identity-service/internal/infrastructure/persistence"
	"identity-service/internal/infrastructure/security"
	uiHttp "identity-service/internal/interfaces/http"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAuthHandler() *uiHttp.AuthHandler {
	ts := testStorage()
	truncateTables(ts.DB)

	userRepo := persistence.NewUserRepository(ts.DB)
	hasher := security.NewPasswordHasher()
	jwt := security.NewJWTService("TestSecretKey")
	refreshTokenRepo := persistence.NewRefreshTokenRepository(ts.Redis)
	refreshTokenService := security.NewRefreshTokenService()

	signInUC := auth.NewSignInUseCase(userRepo, hasher, jwt, refreshTokenRepo, refreshTokenService)

	return uiHttp.NewAuthHandler(signInUC, nil)
}

func TestAuthHandler_SignIn_Success(t *testing.T) {
	ts := testStorage()
	aHandler := setupAuthHandler()
	repo := persistence.NewUserRepository(ts.DB)
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

	repo.Create(context.Background(), newUser)

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
		router.POST("/auth/signin", aHandler.SignIn)

		body, _ := json.Marshal(tc.requestBody)
		req, _ := http.NewRequest(http.MethodPost, "/auth/signin", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		type response struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}

		assert.Equal(t, tc.expectedCode, recorder.Code)

		var res response
		err := json.Unmarshal(recorder.Body.Bytes(), &res)
		assert.NoError(t, err)

		_, err = jwt.ParseToken(res.AccessToken)
		assert.NoError(t, err)

		_, err = hex.DecodeString(res.RefreshToken)
		assert.NoError(t, err)
	})
}

func TestAuthHandler_SignIn_Failed(t *testing.T) {
	ts := testStorage()
	aHandler := setupAuthHandler()
	repo := persistence.NewUserRepository(ts.DB)
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

	repo.Create(context.Background(), newUser)

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
			router.POST("/auth/signin", aHandler.SignIn)

			body, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/auth/signin", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedCode, recorder.Code)
			assert.JSONEq(t, tc.expectedBody, recorder.Body.String())
		})
	}
}
