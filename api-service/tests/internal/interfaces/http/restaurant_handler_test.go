package http_test

import (
	"api-service/internal/application/restaurant"
	dRestaurant "api-service/internal/domain/restaurant"
	"api-service/internal/infrastructure/persistence"
	iValidator "api-service/internal/infrastructure/validator"
	uiHttp "api-service/internal/interfaces/http"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRestaurantHandler() *uiHttp.RestaurantHandler {
	restaurantRepo := persistence.NewRestaurantRepository(testDB)
	createUseCase := restaurant.NewCreateRestaurantUseCase(restaurantRepo)
	customValidator := iValidator.NewCustomValidator(nil, restaurantRepo)

	restaurantUseCases := &uiHttp.RestaurantUseCases{
		CreateRestaurant: createUseCase,
		CustomValidator:  customValidator,
	}

	return uiHttp.NewRestaurantHandler(restaurantUseCases)
}

func TestRestaurantHandler_Create_Success(t *testing.T) {
	rHandler := setupRestaurantHandler()

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
	})
	router.Use(MockAuthMiddleware())
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

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response restaurant.RestaurantResponseDTO
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Test Restaurant", response.Name)
	assert.Equal(t, "test-restaurant", response.Slug)
	assert.Equal(t, "123 Test Street", response.Address)
	assert.Equal(t, uint(1), response.UserID)

	var createdRestaurant dRestaurant.Restaurant
	testDB.First(&createdRestaurant)

	assert.Equal(t, "Test Restaurant", createdRestaurant.Name)
	assert.Equal(t, "test-restaurant", createdRestaurant.Slug)
	assert.Equal(t, "123 Test Street", createdRestaurant.Address)
	assert.Equal(t, uint(1), createdRestaurant.UserID)
}

func TestRestaurantHandler_Create_Failure_ValidationError(t *testing.T) {
	rHandler := setupRestaurantHandler()

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
	})
	router.Use(MockAuthMiddleware())
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

func TestRestaurantHandler_Create_Failure_DuplicateSlug(t *testing.T) {
	rHandler := setupRestaurantHandler()

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
	})
	router.Use(MockAuthMiddleware())
	router.POST("/api/restaurants", rHandler.Create)

	// First request (success)
	validPayload := `{"name": "Pizzeria Uno", "slug": "pizzeria-uno", "address": "123 Main St"}`
	req1, _ := http.NewRequest("POST", "/api/restaurants", bytes.NewBufferString(validPayload))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "Bearer mock-valid-token")

	recorder1 := httptest.NewRecorder()
	router.ServeHTTP(recorder1, req1)
	assert.Equal(t, http.StatusCreated, recorder1.Code)

	// Second request (duplicate slug)
	req2, _ := http.NewRequest("POST", "/api/restaurants", bytes.NewBufferString(validPayload))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer mock-valid-token")

	recorder2 := httptest.NewRecorder()
	router.ServeHTTP(recorder2, req2)
	assert.Equal(t, http.StatusUnprocessableEntity, recorder2.Code)
	assert.Contains(t, recorder2.Body.String(), "slug")
}

func TestRestaurantHandler_Create_Failure_Unauthorized(t *testing.T) {
	rHandler := setupRestaurantHandler()

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
	})
	router.Use(MockAuthMiddleware())
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
