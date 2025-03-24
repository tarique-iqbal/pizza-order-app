package http

import (
	"net/http"

	"pizza-order-api/internal/application/user"
	infrastructureValidator "pizza-order-api/internal/infrastructure/validator"
	"pizza-order-api/internal/interfaces/http/validation"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	useCases *UserUseCases
}

type UserUseCases struct {
	CreateUser      user.CreateUserUseCase
	SignIn          user.SignInUserUseCase
	CustomValidator *infrastructureValidator.CustomValidator
}

func NewUserHandler(useCases *UserUseCases) *UserHandler {
	return &UserHandler{useCases: useCases}
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (h *UserHandler) SignIn(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errors := validation.ExtractValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	token, err := h.useCases.SignIn.Execute(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
