package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"api-service/internal/interfaces/http/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequireRole(t *testing.T) {
	tests := []struct {
		name         string
		role         string
		expectedCode int
	}{
		{"CorrectRole", "Owner", http.StatusOK},
		{"IncorrectRole", "User", http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			r := gin.New()

			r.Use(func(c *gin.Context) {
				c.Set("role", tt.role)
			})
			r.Use(middlewares.RequireRole("Owner"))

			r.GET("/api/restaurants/losteria", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req, _ := http.NewRequest(http.MethodGet, "/api/restaurants/losteria", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
