package http

import (
	"api-service/internal/application/user"
	iValidator "api-service/internal/infrastructure/validator"
	"api-service/internal/interfaces/http/mapper"
	"api-service/internal/interfaces/http/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	useCases *UserUseCases
}

type UserUseCases struct {
	CreateUser      *user.CreateUserUseCase
	CustomValidator *iValidator.CustomValidator
}

func NewUserHandler(useCases *UserUseCases) *UserHandler {
	return &UserHandler{useCases: useCases}
}

func (h *UserHandler) Create(ctx *gin.Context) {
	var input user.UserCreateDTO
	reqCtx := ctx.Request.Context()

	if err := ctx.ShouldBindJSON(&input); err != nil {
		errors := validation.ExtractValidationErrors(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors})
		return
	}

	response, err := h.useCases.CreateUser.Execute(reqCtx, input)
	if err != nil {
		status := mapper.MapErrorToHTTPStatus(err)
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}
