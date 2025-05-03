package http

import (
	"net/http"
	"strconv"

	"api-service/internal/application/restaurant"
	"api-service/internal/interfaces/http/response"
	"api-service/internal/interfaces/http/validation"

	"github.com/gin-gonic/gin"
)

type PizzaSizeHandler struct {
	createPizzaSizeUseCase *restaurant.CreatePizzaSizeUseCase
}

func NewPizzaSizeHandler(createUC *restaurant.CreatePizzaSizeUseCase) *PizzaSizeHandler {
	return &PizzaSizeHandler{
		createPizzaSizeUseCase: createUC,
	}
}

func (h *PizzaSizeHandler) Create(ctx *gin.Context) {
	restaurantIDStr := ctx.Param("id")
	restaurantID, err := strconv.ParseUint(restaurantIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
		return
	}

	var dto restaurant.PizzaSizeCreateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		errors := validation.ExtractValidationErrors(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors})
		return
	}

	ownerID := ctx.MustGet("userID").(uint)
	res, err := h.createPizzaSizeUseCase.Execute(ctx.Request.Context(), uint(restaurantID), ownerID, dto)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, res)
}
