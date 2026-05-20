package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"restaurant-service/internal/interfaces/http/middleware"
)

func TestRequireRole(t *testing.T) {
	tests := []struct {
		name         string
		role         string
		expectedCode int
	}{
		{
			name:         "authorized",
			role:         "owner",
			expectedCode: http.StatusOK,
		},
		{
			name:         "forbidden",
			role:         "user",
			expectedCode: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			r := gin.New()

			r.Use(func(c *gin.Context) {
				c.Set(middleware.CtxUserRole, tt.role)
				c.Next()
			})

			r.Use(middleware.RequireRole("owner"))

			url := "/restaurants/019e26ff-37de-7ed2-b093-b3684613c7cc/address"

			r.PATCH(url, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "success",
				})
			})

			req, err := http.NewRequest(http.MethodPatch, url, nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
