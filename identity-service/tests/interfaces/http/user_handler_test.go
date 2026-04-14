package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	userapp "identity-service/internal/application/user"
	"identity-service/internal/infrastructure/auth"
	"identity-service/internal/infrastructure/persistence"
	"identity-service/internal/infrastructure/security"
	httpui "identity-service/internal/interfaces/http"
	"identity-service/internal/shared/event"
	"identity-service/tests/infrastructure/db/fixtures"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func setupUserHandler() *httpui.UserHandler {
	ts := testStorage()
	truncateTables(ts.DB)

	emailVerificationRepo := persistence.NewEmailVerificationRepository(ts.DB)
	codeVerifier := auth.NewEmailVerifier(emailVerificationRepo)
	userRepo := persistence.NewUserRepository(ts.DB)
	hasher := security.NewPasswordHasher()
	mockPublisher = &MockEventPublisher{}

	register := userapp.NewRegister(codeVerifier, userRepo, hasher, mockPublisher)
	findByID := userapp.NewFindByID(userRepo)

	if err := fixtures.LoadEmailVerificationFixtures(ts.DB); err != nil {
		panic(err)
	}

	return httpui.NewUserHandler(register, findByID)
}

func TestUserHandler_Register_Success(t *testing.T) {
	handler := setupUserHandler()

	router := gin.Default()
	router.POST("/users", handler.Register)

	reqBody, _ := json.Marshal(map[string]string{
		"first_name": "Alice",
		"last_name":  "Doe",
		"email":      "alice@example.com",
		"password":   "pass123",
		"role":       "user",
		"code":       "347578",
	})

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "alice@example.com")
}

func TestUserHandler_FindByID_Success(t *testing.T) {
	ts := testStorage()
	handler := setupUserHandler()

	user, _ := fixtures.CreateUser(ts.DB, "user")

	router := gin.Default()
	router.Use(MockAuthMiddleware(user.ID, user.Role))
	router.GET("/users/:id", handler.FindByID)

	req, _ := http.NewRequest(http.MethodGet, "/users/"+strconv.Itoa(user.ID), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-valid-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var res userapp.Response
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.Equal(t, user.ID, res.ID)
	assert.Equal(t, user.FirstName, res.Name.First)
	assert.Equal(t, user.LastName, res.Name.Last)
	assert.Equal(t, user.Email, res.Email)
	assert.Equal(t, user.Role, res.Role)
	assert.Equal(t, user.Status, res.Status)
}

func TestUserHandler_FindByID_NotFound(t *testing.T) {
	ts := testStorage()
	handler := setupUserHandler()

	user, _ := fixtures.CreateUser(ts.DB, "user")

	router := gin.Default()
	router.Use(MockAuthMiddleware(user.ID, user.Role))
	router.GET("/users/:id", handler.FindByID)

	req := httptest.NewRequest(http.MethodGet, "/users/989", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-valid-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)

	assert.Equal(t, "user not found", body["error"])
}

func TestUserHandler_FindByID_Failure_Unauthorized(t *testing.T) {
	ts := testStorage()
	handler := setupUserHandler()

	user, _ := fixtures.CreateUser(ts.DB, "user")

	router := gin.Default()
	router.Use(MockAuthMiddleware(user.ID, user.Role))
	router.GET("/users/:id", handler.FindByID)

	req, _ := http.NewRequest(http.MethodGet, "/users/"+strconv.Itoa(user.ID), nil)
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")
}
