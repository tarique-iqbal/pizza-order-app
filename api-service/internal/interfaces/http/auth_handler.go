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
	SignIn                  *auth.SignInUseCase
	CreateEmailVerification *auth.CreateEmailVerificationUseCase
}

func NewAuthHandler(useCases *AuthUseCases) *AuthHandler {
	return &AuthHandler{useCases: useCases}
}

func (h *AuthHandler) SignIn(ctx *gin.Context) {
	var dto auth.SignInRequestDTO
	reqCtx := ctx.Request.Context()

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		errors := validation.ExtractValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	token, err := h.useCases.SignIn.Execute(reqCtx, dto.Email, dto.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) CreateEmailVerification(ctx *gin.Context) {
	var dto auth.EmailVerificationRequestDTO
	reqCtx := ctx.Request.Context()

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		errors := validation.ExtractValidationErrors(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors})
		return
	}

	err := h.useCases.CreateEmailVerification.Execute(reqCtx, dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create verification"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
