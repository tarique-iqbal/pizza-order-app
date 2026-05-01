package http_test

import (
	"bytes"
	"context"
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
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var mockPublisher *MockEventPublisher

type MockEventPublisher struct {
	PublishedEvents []event.Event
	PublishedRaw    [][]byte
	ShouldFail      bool
}

func (m *MockEventPublisher) PublishEvent(ctx context.Context, e event.Event) error {
	m.PublishedEvents = append(m.PublishedEvents, e)
	if m.ShouldFail {
		return errors.New("mock publish failure")
	}
	return nil
}

func (m *MockEventPublisher) PublishRaw(ctx context.Context, topic string, jsonData []byte) error {
	m.PublishedRaw = append(m.PublishedRaw, jsonData)
	if m.ShouldFail {
		return errors.New("mock raw publish failure")
	}
	return nil
}

func setupUserHandler() *httpui.UserHandler {
	ts := testStorage()
	truncateTables(ts.DB)

	emailVerificationRepo := persistence.NewEmailVerificationRepository(ts.DB)
	codeVerifier := auth.NewEmailVerifier(emailVerificationRepo)
	hasher := security.NewPasswordHasher()
	userRepo := persistence.NewUserRepository(ts.DB)
	outboxRepo := persistence.NewOutboxRepository(ts.DB)
	mockPublisher = &MockEventPublisher{}

	register := userapp.NewRegisterCustomer(codeVerifier, userRepo, hasher, mockPublisher)
	registerOwner := userapp.NewRegisterOwner(ts.DB, codeVerifier, hasher, userRepo, outboxRepo, mockPublisher)
	findByID := userapp.NewFindByID(userRepo)

	_ = fixtures.LoadEmailVerificationFixtures(ts.DB)

	return httpui.NewUserHandler(register, registerOwner, findByID)
}

func TestUserHandler_RegisterOwner(t *testing.T) {
	handler := setupUserHandler()

	tests := []struct {
		name           string
		body           map[string]string
		rawBody        string
		expectedStatus int
		expectError    bool
	}{
		{
			name: "success",
			body: map[string]string{
				"firstName":    "Alice",
				"lastName":     "Doe",
				"email":        "alice@example.com",
				"password":     "pass123",
				"code":         "347578", // from fixture
				"businessName": "Domino's Pizza",
				"vatNumber":    "DE123456789",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid json",
			rawBody:        "invalid-json",
			expectedStatus: http.StatusUnprocessableEntity,
			expectError:    true,
		},
		{
			name: "invalid code",
			body: map[string]string{
				"firstName":    "Alice",
				"lastName":     "Doe",
				"email":        "alice@example.com",
				"password":     "pass123",
				"code":         "000000", // invalid
				"businessName": "Domino's Pizza",
				"vatNumber":    "DE987654321",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := testStorage()
			truncateTables(ts.DB)

			_ = fixtures.LoadEmailVerificationFixtures(ts.DB)
			_ = fixtures.LoadUserFixtures(ts.DB)

			router := gin.Default()
			router.POST("/owners", handler.RegisterOwner)

			var reqBody []byte
			if tt.rawBody != "" {
				reqBody = []byte(tt.rawBody)
			} else {
				reqBody, _ = json.Marshal(tt.body)
			}

			req, _ := http.NewRequest("POST", "/owners", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectError {
				assert.Contains(t, w.Body.String(), "error")
			}
		})
	}
}

func TestUserHandler_RegisterCustomer(t *testing.T) {
	handler := setupUserHandler()

	tests := []struct {
		name           string
		body           map[string]string
		rawBody        string
		expectedStatus int
		expectError    bool
	}{
		{
			name: "success",
			body: map[string]string{
				"firstName": "Sophie",
				"lastName":  "Mueller",
				"email":     "sophie.mueller@example.com",
				"password":  "pass123",
				"code":      "365189", // from fixture
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "validation error",
			body:           map[string]string{},
			expectedStatus: http.StatusUnprocessableEntity,
			expectError:    true,
		},
		{
			name: "duplicate email",
			body: map[string]string{
				"firstName": "Existing",
				"lastName":  "User",
				"email":     "existing@example.com", // from fixture
				"password":  "pass123",
				"code":      "347578",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := testStorage()
			truncateTables(ts.DB)

			_ = fixtures.LoadEmailVerificationFixtures(ts.DB)
			_ = fixtures.LoadUserFixtures(ts.DB)

			router := gin.Default()
			router.POST("/customers", handler.RegisterCustomer)

			var reqBody []byte
			if tt.rawBody != "" {
				reqBody = []byte(tt.rawBody)
			} else {
				reqBody, _ = json.Marshal(tt.body)
			}

			req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectError {
				assert.Contains(t, w.Body.String(), "error")
			}
		})
	}
}

func TestUserHandler_FindByID_Success(t *testing.T) {
	ts := testStorage()
	handler := setupUserHandler()

	u := fixtures.NewUser()
	user, err := fixtures.CreateUser(ts.DB, u)
	require.NoError(t, err)

	router := gin.Default()
	router.Use(MockAuthMiddleware(user.ID.String(), user.Role))
	router.GET("/users/:id", handler.FindByID)

	req, _ := http.NewRequest(http.MethodGet, "/users/"+user.ID.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-valid-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var res userapp.Response
	err = json.Unmarshal(w.Body.Bytes(), &res)
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

	u := fixtures.NewUser()
	user, err := fixtures.CreateUser(ts.DB, u)
	require.NoError(t, err)

	router := gin.Default()
	router.Use(MockAuthMiddleware(user.ID.String(), user.Role))
	router.GET("/users/:id", handler.FindByID)

	newID, _ := uuid.NewV7()

	req := httptest.NewRequest(http.MethodGet, "/users/"+newID.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-valid-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var body map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)

	assert.Equal(t, "user not found", body["error"])
}

func TestUserHandler_FindByID_Failure_Unauthorized(t *testing.T) {
	ts := testStorage()
	handler := setupUserHandler()

	u := fixtures.NewUser()
	user, err := fixtures.CreateUser(ts.DB, u)
	require.NoError(t, err)

	router := gin.Default()
	router.Use(MockAuthMiddleware(user.ID.String(), user.Role))
	router.GET("/users/:id", handler.FindByID)

	req, _ := http.NewRequest(http.MethodGet, "/users/"+user.ID.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")
}
