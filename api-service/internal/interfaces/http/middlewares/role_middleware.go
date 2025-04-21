package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(expectedRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.MustGet("role").(string)
		if role != expectedRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		c.Next()
	}
}
