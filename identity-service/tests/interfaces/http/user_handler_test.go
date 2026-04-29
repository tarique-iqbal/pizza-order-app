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
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	hasher := security.NewPasswordHasher()
	userRepo := persistence.NewUserRepository(ts.DB)
	outboxRepo := persistence.NewOutboxRepository(ts.DB)
	mockPublisher = &MockEventPublisher{}

	register := userapp.NewRegister(codeVerifier, userRepo, hasher, mockPublisher)
	registerOwner := userapp.NewRegisterOwner(ts.DB, codeVerifier, hasher, userRepo, outboxRepo, mockPublisher)
	findByID := userapp.NewFindByID(userRepo)

	if err := fixtures.LoadEmailVerificationFixtures(ts.DB); err != nil {
		panic(err)
	}

	return httpui.NewUserHandler(register, registerOwner, findByID)
}

func TestUserHandler_Register(t *testing.T) {
	handler := setupUserHandler()

	tests := []struct {
		name           string
		role           string
		body           map[string]string
		rawBody        string
		expectedStatus int
		expectError    bool
	}{
		{
			name: "owner success",
			body: map[string]string{
				"firstName":    "Alice",
				"lastName":     "Doe",
				"email":        "alice@example.com",
				"password":     "pass123",
				"code":         "347578", // fixture
				"role":         "owner",
				"businessName": "Domino's Pizza",
				"vatNumber":    "DE987654321",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "owner validation error",
			role:           "owner",
			rawBody:        "invalid-json",
			expectedStatus: http.StatusUnprocessableEntity,
			expectError:    true,
		},
		{
			name: "owner invalid code",
			body: map[string]string{
				"firstName":    "Alice",
				"lastName":     "Doe",
				"email":        "alice@example.com",
				"password":     "pass123",
				"code":         "000000", // invalid
				"role":         "owner",
				"businessName": "Domino's Pizza",
				"vatNumber":    "DE987654321",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "user success",
			body: map[string]string{
				"firstName": "Sophie",
				"lastName":  "Mueller",
				"email":     "sophie.mueller@example.com",
				"password":  "pass123",
				"code":      "365189", // fixture
				"role":      "user",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "user validation error",
			body:           map[string]string{},
			expectedStatus: http.StatusUnprocessableEntity,
			expectError:    true,
		},
		{
			name: "user duplicate email",
			body: map[string]string{
				"firstName": "Existing",
				"lastName":  "User",
				"email":     "existing@example.com", // fixture
				"password":  "pass123",
				"code":      "347578",
				"role":      "user",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := testStorage()
			truncateTables(ts.DB)

			if err := fixtures.LoadEmailVerificationFixtures(ts.DB); err != nil {
				panic(err)
			}

			if err := fixtures.LoadUserFixtures(ts.DB); err != nil {
				panic(err)
			}

			router := gin.Default()
			router.POST("/users", handler.Register)

			var reqBody []byte
			if tt.rawBody != "" {
				reqBody = []byte(tt.rawBody)
			} else {
				reqBody, _ = json.Marshal(tt.body)
			}

			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
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
