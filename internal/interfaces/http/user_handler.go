package http

import (
	"net/http"
	"pizza-order-api/internal/application/user"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	createUserUseCase *user.CreateUserUseCase
}

func NewUserHandler(createUserUseCase *user.CreateUserUseCase) *UserHandler {
	return &UserHandler{createUserUseCase}
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var input user.UserCreateDTO

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.createUserUseCase.Execute(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}
