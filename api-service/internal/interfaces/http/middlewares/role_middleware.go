package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(expectedRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := ctx.MustGet("role").(string)
		if role != expectedRole {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		ctx.Next()
	}
}
