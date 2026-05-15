package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"restaurant-service/internal/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequireRole(t *testing.T) {
	tests := []struct {
		name         string
		role         string
		expectedCode int
	}{
		{"CorrectRole", "owner", http.StatusOK},
		{"IncorrectRole", "user", http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			r := gin.New()

			r.Use(func(c *gin.Context) {
				c.Set("userRole", tt.role)
				c.Next()
			})

			r.Use(middleware.RequireRole("owner"))

			url := "/restaurants/019e26ff-37de-7ed2-b093-b3684613c7cc/addresses"

			r.GET(url, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "success",
				})
			})

			req, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
