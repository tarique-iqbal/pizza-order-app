package http

import (
	"api-service/internal/application/auth"
	"api-service/internal/interfaces/http/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	useCases *AuthUseCases
}

type AuthUseCases struct {
	SignIn *auth.SignInUseCase
}

func NewAuthHandler(useCases *AuthUseCases) *AuthHandler {
	return &AuthHandler{useCases: useCases}
}

func (h *AuthHandler) SignIn(ctx *gin.Context) {
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
