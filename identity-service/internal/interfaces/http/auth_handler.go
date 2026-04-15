package http

import (
	"identity-service/internal/application/auth"
	"identity-service/internal/interfaces/http/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	login        *auth.Login
	emailOTP     *auth.RequestEmailOTP
	refreshToken *auth.RefreshToken
	logout       *auth.Logout
}

func NewAuthHandler(
	login *auth.Login,
	emailOTP *auth.RequestEmailOTP,
	refreshToken *auth.RefreshToken,
	logout *auth.Logout,
) *AuthHandler {
	return &AuthHandler{
		login:        login,
		emailOTP:     emailOTP,
		refreshToken: refreshToken,
		logout:       logout,
	}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var input auth.LoginRequest
	reqCtx := ctx.Request.Context()

	if err := ctx.ShouldBindJSON(&input); err != nil {
		errors := validation.ExtractValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	response, err := h.login.Execute(reqCtx, input)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *AuthHandler) CreateEmailVerification(ctx *gin.Context) {
	var input auth.EmailVerificationRequest
	reqCtx := ctx.Request.Context()

	if err := ctx.ShouldBindJSON(&input); err != nil {
		errors := validation.ExtractValidationErrors(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors})
		return
	}

	err := h.emailOTP.Execute(reqCtx, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create verification"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *AuthHandler) Refresh(ctx *gin.Context) {
	var input auth.RefreshRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	response, err := h.refreshToken.Execute(ctx, input)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	var input auth.LogoutRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if _, exists := ctx.Get("userID"); !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.logout.Execute(ctx, input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "logged out successfully",
	})
}
