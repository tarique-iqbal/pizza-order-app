package handlers_test

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"

	"restaurant-service/internal/interfaces/http/middleware"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func MockAuthMiddleware(userID, role string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(middleware.CtxUserID, userID)
		ctx.Set(middleware.CtxUserRole, role)
		ctx.Next()
	}
}
