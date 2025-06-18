package http_test

import (
	"api-service/internal/application/restaurant"
	dRestaurant "api-service/internal/domain/restaurant"
	"api-service/internal/infrastructure/persistence"
	uiHttp "api-service/internal/interfaces/http"
	"api-service/internal/interfaces/http/middlewares"
	"api-service/tests/infrastructure/db/fixtures"
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
	restAddrRepo := persistence.NewRestaurantAddressRepository(testDB)
	createRestaurantUC := restaurant.NewCreateRestaurantUseCase(testDB, restaurantRepo, restAddrRepo)

	return uiHttp.NewRestaurantHandler(createRestaurantUC)
}

func TestRestaurantHandler_Create_Success(t *testing.T) {
	rHandler := setupRestaurantHandler(t)
	usr, _ := fixtures.CreateUser(testDB, "owner")

	router := gin.Default()
	router.Use(MockAuthMiddleware(usr.ID, usr.Role), middlewares.RequireRole("Owner"))
	router.POST("/api/restaurants", rHandler.Create)

	reqBody := map[string]interface{}{
		"name":          "Test Restaurant",
		"email":         "unique@test.com",
		"phone":         "+49 89 22334455",
		"house":         "1",
		"street":        "Main Str.",
		"city":          "Cityville",
		"postal_code":   "12345",
		"delivery_type": "own_delivery",
		"delivery_km":   5,
		"specialties":   []string{"italian", "wood_fired"},
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
	assert.Equal(t, "unique@test.com", response.Email)
	assert.Equal(t, "test-restaurant-cityville", response.Slug)
	assert.Equal(t, usr.ID, response.UserID)

	var createdRestaurant dRestaurant.Restaurant
	testDB.Where("slug = ?", "test-restaurant-cityville").First(&createdRestaurant)

	assert.Equal(t, "Test Restaurant", createdRestaurant.Name)
	assert.Equal(t, "unique@test.com", createdRestaurant.Email)
	assert.Equal(t, usr.ID, createdRestaurant.UserID)
}

func TestRestaurantHandler_Create_Failure_ValidationError(t *testing.T) {
	rHandler := setupRestaurantHandler(t)
	usr, _ := fixtures.CreateUser(testDB, "owner")

	router := gin.Default()
	router.Use(MockAuthMiddleware(usr.ID, usr.Role), middlewares.RequireRole("Owner"))
	router.POST("/api/restaurants", rHandler.Create)

	payload := `{"name": "Pizza Restaurant"}`

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
	usr, _ := fixtures.CreateUser(testDB, "owner")

	router := gin.Default()
	router.Use(MockAuthMiddleware(usr.ID, usr.Role), middlewares.RequireRole("Owner"))
	router.POST("/api/restaurants", rHandler.Create)

	reqBody := map[string]interface{}{
		"name":          "Test Restaurant",
		"email":         "unique@test.com",
		"phone":         "+49 89 22334455",
		"house":         "1",
		"street":        "Main Str.",
		"city":          "Cityville",
		"postal_code":   "12345",
		"delivery_type": "own_delivery",
		"delivery_km":   5,
		"specialties":   []string{"italian", "wood_fired"},
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/restaurants", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Unauthorized")
}

func TestRestaurantHandler_Create_Failure_UnauthorizedRole(t *testing.T) {
	rHandler := setupRestaurantHandler(t)
	usr, _ := fixtures.CreateUser(testDB, "user")

	router := gin.Default()
	router.Use(MockAuthMiddleware(usr.ID, usr.Role), middlewares.RequireRole("Owner"))
	router.POST("/api/restaurants", rHandler.Create)

	reqBody := map[string]interface{}{
		"name":          "Test Restaurant",
		"email":         "unique@test.com",
		"phone":         "+49 89 22334455",
		"house":         "1",
		"street":        "Main Str.",
		"city":          "Cityville",
		"postal_code":   "12345",
		"delivery_type": "own_delivery",
		"delivery_km":   5,
		"specialties":   []string{"italian", "wood_fired"},
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
