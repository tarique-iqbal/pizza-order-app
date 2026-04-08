package http

import (
	"identity-service/internal/application/auth"
	"identity-service/internal/interfaces/http/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	login    *auth.Login
	emailOTP *auth.RequestEmailOTP
}

func NewAuthHandler(
	login *auth.Login,
	emailOTP *auth.RequestEmailOTP,
) *AuthHandler {
	return &AuthHandler{login: login, emailOTP: emailOTP}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var dto auth.LoginRequest
	reqCtx := ctx.Request.Context()

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		errors := validation.ExtractValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	response, err := h.login.Execute(reqCtx, dto.Email, dto.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *AuthHandler) CreateEmailVerification(ctx *gin.Context) {
	var dto auth.EmailVerificationRequest
	reqCtx := ctx.Request.Context()

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		errors := validation.ExtractValidationErrors(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors})
		return
	}

	err := h.emailOTP.Execute(reqCtx, dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create verification"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
