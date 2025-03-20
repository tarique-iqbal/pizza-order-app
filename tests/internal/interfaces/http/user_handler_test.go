package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	appUser "pizza-order-api/internal/application/user"
	domainUser "pizza-order-api/internal/domain/user"
	"pizza-order-api/internal/infrastructure/persistence"
	"pizza-order-api/internal/infrastructure/security"
	interfacesHttp "pizza-order-api/internal/interfaces/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&domainUser.User{})
	return db
}

func TestCreateUserHandler(t *testing.T) {
	db := setupTestDB()
	userRepo := persistence.NewUserRepository(db)
	createUserUseCase := appUser.NewCreateUserUseCase(userRepo)
	handler := interfacesHttp.NewUserHandler(createUserUseCase, nil)

	router := gin.Default()
	router.POST("/api/users/signup", handler.CreateUser)

	reqBody, _ := json.Marshal(map[string]string{
		"first_name": "Alice",
		"last_name":  "Doe",
		"email":      "alice@example.com",
		"password":   "pass123",
		"role":       "user",
	})

	req, _ := http.NewRequest("POST", "/api/users/signup", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestUserHandler_SignIn_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB()
	userRepo := persistence.NewUserRepository(db)
	signInUserUseCase := appUser.NewSignInUserUseCase(userRepo)
	userHandler := interfacesHttp.NewUserHandler(nil, signInUserUseCase)

	repo := persistence.NewUserRepository(db)
	hp, _ := security.HashPassword("password123")

	newUser := &domainUser.User{
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
		router.POST("/api/users/signin", userHandler.SignIn)

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
	gin.SetMode(gin.TestMode)

	db := setupTestDB()
	userRepo := persistence.NewUserRepository(db)
	signInUserUseCase := appUser.NewSignInUserUseCase(userRepo)
	userHandler := interfacesHttp.NewUserHandler(nil, signInUserUseCase)

	repo := persistence.NewUserRepository(db)
	hp, _ := security.HashPassword("password123")

	newUser := &domainUser.User{
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
			expectedBody: `{"error":"Key: 'Email' Error:Field validation for 'Email' failed on the 'email' tag"}`,
		},
		{
			name: "Invalid Request Body",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Key: 'Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.Default()
			router.POST("/api/users/signin", userHandler.SignIn)

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
