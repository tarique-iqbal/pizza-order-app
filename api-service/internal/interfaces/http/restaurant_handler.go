package http

import (
	"api-service/internal/application/restaurant"
	"api-service/internal/interfaces/http/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RestaurantHandler struct {
	createRestaurantUC *restaurant.CreateRestaurantUseCase
}

func NewRestaurantHandler(
	createRestaurantUC *restaurant.CreateRestaurantUseCase,
) *RestaurantHandler {
	return &RestaurantHandler{createRestaurantUC: createRestaurantUC}
}

func (h *RestaurantHandler) Create(ctx *gin.Context) {
	var dto restaurant.RestaurantCreateDTO
	reqCtx := ctx.Request.Context()

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		errors := validation.ExtractValidationErrors(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors})
		return
	}

	dto.UserID = ctx.MustGet("userID").(uint)

	res, err := h.createRestaurantUC.Execute(reqCtx, dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create restaurant"})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}
