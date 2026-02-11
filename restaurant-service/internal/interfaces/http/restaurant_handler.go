package http

import (
	"net/http"
	"restaurant-service/internal/application/restaurant"
	"restaurant-service/internal/interfaces/http/response"
	"restaurant-service/internal/interfaces/http/validation"

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
		response.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, res)
}
