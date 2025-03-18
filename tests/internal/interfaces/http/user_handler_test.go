package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	appUser "pizza-order-api/internal/application/user"
	domainUser "pizza-order-api/internal/domain/user"
	"pizza-order-api/internal/infrastructure/persistence"
	interfacesHttp "pizza-order-api/internal/interfaces/http"
	"testing"

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
	useCase := appUser.NewCreateUserUseCase(userRepo)
	handler := interfacesHttp.NewUserHandler(useCase)

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
