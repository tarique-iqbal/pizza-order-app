package http

import (
	"identity-service/internal/application/user"
	"identity-service/internal/interfaces/http/mapper"
	"identity-service/internal/interfaces/http/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	register *user.Register
}

func NewUserHandler(reg *user.Register) *UserHandler {
	return &UserHandler{register: reg}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var input user.RegisterRequest
	reqCtx := ctx.Request.Context()

	if err := ctx.ShouldBindJSON(&input); err != nil {
		errors := validation.ExtractValidationErrors(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors})
		return
	}

	response, err := h.register.Execute(reqCtx, input)
	if err != nil {
		status := mapper.MapErrorToHTTPStatus(err)
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}
