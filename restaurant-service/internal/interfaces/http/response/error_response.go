package response

import (
	"net/http"

	"restaurant-service/internal/domain/restaurant"
	apperr "restaurant-service/internal/shared/errors"

	"github.com/gin-gonic/gin"
)

func HandleError(ctx *gin.Context, err error) {
	switch err {
	case apperr.ErrUnauthorized:
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case apperr.ErrForbidden:
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case apperr.ErrNotFound:
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case apperr.ErrConflict,
		restaurant.ErrPizzaSizeAlreadyExists,
		restaurant.ErrEmailAlreadyExists:
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
