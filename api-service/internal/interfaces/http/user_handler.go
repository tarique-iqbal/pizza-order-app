package http

import (
	"api-service/internal/application/user"
	"api-service/internal/interfaces/http/mapper"
	"api-service/internal/interfaces/http/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	createUserUC *user.CreateUserUseCase
}

func NewUserHandler(createUserUC *user.CreateUserUseCase) *UserHandler {
	return &UserHandler{createUserUC: createUserUC}
}

func (h *UserHandler) Create(ctx *gin.Context) {
	var input user.UserCreateDTO
	reqCtx := ctx.Request.Context()

	if err := ctx.ShouldBindJSON(&input); err != nil {
		errors := validation.ExtractValidationErrors(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors})
		return
	}

	response, err := h.createUserUC.Execute(reqCtx, input)
	if err != nil {
		status := mapper.MapErrorToHTTPStatus(err)
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}
