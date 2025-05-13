package http_test

import (
	"api-service/internal/application/restaurant"
	dRestaurant "api-service/internal/domain/restaurant"
	"api-service/internal/infrastructure/persistence"
	uiHttp "api-service/internal/interfaces/http"
	"api-service/internal/interfaces/http/middlewares"
	"api-service/tests/internal/infrastructure/db/fixtures"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRestaurantHandler(t *testing.T) *uiHttp.RestaurantHandler {
	resetTables(t)

	restaurantRepo := persistence.NewRestaurantRepository(testDB)
	createUseCase := restaurant.NewCreateRestaurantUseCase(restaurantRepo)

	restaurantUseCases := &uiHttp.RestaurantUseCases{
		CreateRestaurant: createUseCase,
	}

	return uiHttp.NewRestaurantHandler(restaurantUseCases)
}

func TestRestaurantHandler_Create_Success(t *testing.T) {
	rHandler := setupRestaurantHandler(t)
	usr, _ := fixtures.CreateUser(testDB, "Owner")

	router := gin.Default()
	router.Use(MockAuthMiddleware(usr.ID, usr.Role), middlewares.RequireRole("Owner"))
	router.POST("/api/restaurants", rHandler.Create)

	reqBody := map[string]string{
		"name":    "Test Restaurant",
		"slug":    "test-restaurant",
		"address": "123 Test Street",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/restaurants", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-valid-token")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	var response restaurant.RestaurantResponseDTO
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "Test Restaurant", response.Name)
	assert.Equal(t, "test-restaurant", response.Slug)
	assert.Equal(t, "123 Test Street", response.Address)
	assert.Equal(t, usr.ID, response.UserID)

	var createdRestaurant dRestaurant.Restaurant
	testDB.Where("slug = ?", "test-restaurant").First(&createdRestaurant)

	assert.Equal(t, "Test Restaurant", createdRestaurant.Name)
	assert.Equal(t, "test-restaurant", createdRestaurant.Slug)
	assert.Equal(t, "123 Test Street", createdRestaurant.Address)
	assert.Equal(t, usr.ID, createdRestaurant.UserID)
}

func TestRestaurantHandler_Create_Failure_ValidationError(t *testing.T) {
	rHandler := setupRestaurantHandler(t)
	usr, _ := fixtures.CreateUser(testDB, "Owner")

	router := gin.Default()
	router.Use(MockAuthMiddleware(usr.ID, usr.Role), middlewares.RequireRole("Owner"))
	router.POST("/api/restaurants", rHandler.Create)

	payload := `{"slug": "valid-slug"}`

	req, _ := http.NewRequest("POST", "/api/restaurants", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-valid-token")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "errors")
}

func TestRestaurantHandler_Create_Failure_Unauthorized(t *testing.T) {
	rHandler := setupRestaurantHandler(t)
	usr, _ := fixtures.CreateUser(testDB, "Owner")

	router := gin.Default()
	router.Use(MockAuthMiddleware(usr.ID, usr.Role), middlewares.RequireRole("Owner"))
	router.POST("/api/restaurants", rHandler.Create)

	validPayload := `{"name": "New Restaurant", "slug": "new-restaurant", "address": "456 Elm St"}`

	req, _ := http.NewRequest("POST", "/api/restaurants", bytes.NewBufferString(validPayload))
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Unauthorized")
}

func TestRestaurantHandler_Create_Failure_UnauthorizedRole(t *testing.T) {
	rHandler := setupRestaurantHandler(t)
	usr, _ := fixtures.CreateUser(testDB, "User")

	router := gin.Default()
	router.Use(MockAuthMiddleware(usr.ID, usr.Role), middlewares.RequireRole("Owner"))
	router.POST("/api/restaurants", rHandler.Create)

	reqBody := map[string]string{
		"name":    "Test Restaurant",
		"slug":    "test-restaurant",
		"address": "123 Test Street",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/restaurants", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-valid-token")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Access denied")
}
