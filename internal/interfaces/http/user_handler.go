package http

import (
	"net/http"
	"pizza-order-api/internal/application/user"
	"pizza-order-api/internal/interfaces/http/validation"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	createUserUseCase user.CreateUserUseCase
	signInUseCase     user.SignInUserUseCase
}

func NewUserHandler(createUserUseCase user.CreateUserUseCase, signInUseCase user.SignInUserUseCase) *UserHandler {
	return &UserHandler{
		createUserUseCase: createUserUseCase,
		signInUseCase:     signInUseCase,
	}
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var input user.UserCreateDTO

	if err := ctx.ShouldBindJSON(&input); err != nil {
		errors := validation.ExtractValidationErrors(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors})
		return
	}

	response, err := h.createUserUseCase.Execute(input)
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

	token, err := h.signInUseCase.Execute(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
