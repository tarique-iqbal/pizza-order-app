package response

import (
	"net/http"

	dRestaurant "restaurant-service/internal/domain/restaurant"
	sErrors "restaurant-service/internal/shared/errors"

	"github.com/gin-gonic/gin"
)

func HandleError(ctx *gin.Context, err error) {
	switch err {
	case sErrors.ErrUnauthorized:
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case sErrors.ErrForbidden:
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case sErrors.ErrNotFound:
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case sErrors.ErrConflict,
		dRestaurant.ErrPizzaSizeAlreadyExists,
		dRestaurant.ErrEmailAlreadyExists:
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
