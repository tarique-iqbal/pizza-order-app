package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.GetHeader("X-User-ID")
		userRole := ctx.GetHeader("X-User-Role")

		if userID == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing X-User-ID header",
			})
			return
		}

		if userRole == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing X-User-Role header",
			})
			return
		}

		ctx.Set("userID", userID)
		ctx.Set("userRole", userRole)

		ctx.Next()
	}
}
