package http

import (
	"api-service/internal/application/user"
	iValidator "api-service/internal/infrastructure/validator"
	"api-service/internal/interfaces/http/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("uniqueEmail", h.useCases.CustomValidator.UniqueEmail)
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		errors := validation.ExtractValidationErrors(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors})
		return
	}

	response, err := h.useCases.CreateUser.Execute(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}
