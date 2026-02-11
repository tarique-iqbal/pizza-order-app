package http_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"restaurant-service/internal/application/restaurant"
	"restaurant-service/internal/infrastructure/persistence"
	uiHttp "restaurant-service/internal/interfaces/http"
	"restaurant-service/internal/interfaces/http/middlewares"
	"restaurant-service/tests/infrastructure/db/fixtures"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupPizzaSizeHandler(t *testing.T) *uiHttp.PizzaSizeHandler {
	resetTables(t)

	pizzaSizeRepo := persistence.NewPizzaSizeRepository(testDB)
	restaurantRepo := persistence.NewRestaurantRepository(testDB)
	createPizzaSizeUC := restaurant.NewCreatePizzaSizeUseCase(pizzaSizeRepo, restaurantRepo)

	return uiHttp.NewPizzaSizeHandler(createPizzaSizeUC)
}

func TestPizzaSizeHandler_Create_Success(t *testing.T) {
	psHandler := setupPizzaSizeHandler(t)
	rest, _ := fixtures.CreateRestaurant(testDB)

	router := gin.Default()
	router.Use(MockAuthMiddleware(rest.UserID, "owner"), middlewares.RequireRole("owner"))
	router.POST("/restaurants/:id/pizza-sizes", psHandler.Create)

	body := map[string]interface{}{
		"title": "Large",
		"size":  14,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/restaurants/%d/pizza-sizes", rest.ID), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-valid-token")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	var response restaurant.PizzaSizeResponseDTO
	_ = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, "Large", response.Title)
	assert.Equal(t, 14, response.Size)
	assert.Equal(t, rest.ID, response.RestaurantID)
}

func TestPizzaSizeHandler_Create_Failure_ValidationError(t *testing.T) {
	psHandler := setupPizzaSizeHandler(t)
	rest, _ := fixtures.CreateRestaurant(testDB)

	router := gin.Default()
	router.Use(MockAuthMiddleware(rest.UserID, "owner"), middlewares.RequireRole("owner"))
	router.POST("/restaurants/:id/pizza-sizes", psHandler.Create)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/restaurants/%d/pizza-sizes", rest.ID), bytes.NewBufferString(`invalid_json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-valid-token")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "errors")
}

func TestPizzaSizeHandler_Create_Invalid_Restaurant_Id_Param(t *testing.T) {
	psHandler := setupPizzaSizeHandler(t)
	rest, _ := fixtures.CreateRestaurant(testDB)

	router := gin.Default()
	router.Use(MockAuthMiddleware(rest.UserID, "owner"), middlewares.RequireRole("owner"))
	router.POST("/restaurants/:id/pizza-sizes", psHandler.Create)

	req, _ := http.NewRequest("POST", "/restaurants/invalid-id/pizza-sizes", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-valid-token")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Invalid restaurant ID")
}

func TestPizzaSizeHandler_Create_Forbidden_User(t *testing.T) {
	notOwnerID := uint(99)
	psHandler := setupPizzaSizeHandler(t)
	rest, _ := fixtures.CreateRestaurant(testDB)

	router := gin.Default()
	router.Use(MockAuthMiddleware(notOwnerID, "owner"), middlewares.RequireRole("owner"))
	router.POST("/restaurants/:id/pizza-sizes", psHandler.Create)

	body := map[string]interface{}{
		"title": "Unauthorized Pizza",
		"size":  12,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/restaurants/%d/pizza-sizes", rest.ID), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock-valid-token")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "forbidden")
}
