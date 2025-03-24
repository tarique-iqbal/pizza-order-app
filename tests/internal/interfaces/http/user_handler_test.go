package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	applicationUser "pizza-order-api/internal/application/user"
	domainUser "pizza-order-api/internal/domain/user"
	"pizza-order-api/internal/infrastructure/persistence"
	"pizza-order-api/internal/infrastructure/security"
	infrastructureValidator "pizza-order-api/internal/infrastructure/validator"
	interfacesHttp "pizza-order-api/internal/interfaces/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var userHandler *interfacesHttp.UserHandler

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&domainUser.User{})

	return db
}

func setupUserHandler() *interfacesHttp.UserHandler {
	userRepo := persistence.NewUserRepository(testDB)
	createUserUseCase := applicationUser.NewCreateUserUseCase(userRepo)
	signInUserUseCase := applicationUser.NewSignInUserUseCase(userRepo)
	customValidator := infrastructureValidator.NewCustomValidator(userRepo)
	userUseCases := &interfacesHttp.UserUseCases{
		CreateUser:      createUserUseCase,
		SignIn:          signInUserUseCase,
		CustomValidator: customValidator,
	}

	return interfacesHttp.NewUserHandler(userUseCases)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	testDB = setupTestDB()
	userHandler = setupUserHandler()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestUserHandler_CreateUser_Success(t *testing.T) {
	router := gin.Default()
	router.POST("/api/users/signup", userHandler.CreateUser)

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
	repo := persistence.NewUserRepository(testDB)
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
	repo := persistence.NewUserRepository(testDB)
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
