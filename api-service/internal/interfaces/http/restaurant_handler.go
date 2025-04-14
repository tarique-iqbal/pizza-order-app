package http

import (
	"api-service/internal/application/restaurant"
	iValidator "api-service/internal/infrastructure/validator"
	"api-service/internal/interfaces/http/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type RestaurantHandler struct {
	useCases *RestaurantUseCases
}

type RestaurantUseCases struct {
	CreateRestaurant *restaurant.CreateRestaurantUseCase
	CustomValidator  *iValidator.CustomValidator
}

func NewRestaurantHandler(useCases *RestaurantUseCases) *RestaurantHandler {
	return &RestaurantHandler{useCases: useCases}
}

func (h *RestaurantHandler) Create(ctx *gin.Context) {
	var dto restaurant.RestaurantCreateDTO

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("uniqueRSlug", h.useCases.CustomValidator.UniqueRestaurantSlug)
	}

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		errors := validation.ExtractValidationErrors(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors})
		return
	}

	dto.UserID = ctx.MustGet("userID").(uint)

	res, err := h.useCases.CreateRestaurant.Execute(dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create restaurant"})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}
