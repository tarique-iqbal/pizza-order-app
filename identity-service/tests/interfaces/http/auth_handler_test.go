package http_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	authapp "identity-service/internal/application/auth"
	"identity-service/internal/domain/auth"
	"identity-service/internal/domain/user"
	"identity-service/internal/infrastructure/persistence"
	"identity-service/internal/infrastructure/security"
	httpui "identity-service/internal/interfaces/http"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAuthHandler() (
	*httpui.AuthHandler,
	user.UserRepository,
	auth.PasswordHasher,
	auth.JWTManager,
) {
	ts := testStorage()
	truncateTables(ts.DB)

	userRepo := persistence.NewUserRepository(ts.DB)
	hasher := security.NewPasswordHasher()
	jwt := security.NewJWTManager("TestSecretKey")
	refreshTokenRepo := persistence.NewRefreshTokenRepository(ts.Redis)
	refreshTokenManager := security.NewRefreshTokenManager()

	login := authapp.NewLogin(userRepo, hasher, jwt, refreshTokenRepo, refreshTokenManager)
	refreshToken := authapp.NewRefreshToken(jwt, refreshTokenRepo, refreshTokenManager)

	handler := httpui.NewAuthHandler(login, nil, refreshToken)

	return handler, userRepo, hasher, jwt
}

func TestAuthHandler_Login_Success(t *testing.T) {
	handler, repo, hasher, jwt := setupAuthHandler()
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
		router.POST("/auth/login", handler.Login)

		body, _ := json.Marshal(tc.requestBody)
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		type response struct {
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
		}

		assert.Equal(t, tc.expectedCode, recorder.Code)

		var res response
		err := json.Unmarshal(recorder.Body.Bytes(), &res)
		assert.NoError(t, err)

		_, err = jwt.Parse(res.AccessToken)
		assert.NoError(t, err)

		_, err = hex.DecodeString(res.RefreshToken)
		assert.NoError(t, err)
	})
}

func TestAuthHandler_Login_Failed(t *testing.T) {
	handler, repo, hasher, _ := setupAuthHandler()
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
			router.POST("/auth/login", handler.Login)

			body, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedCode, recorder.Code)
			assert.JSONEq(t, tc.expectedBody, recorder.Body.String())
		})
	}
}
