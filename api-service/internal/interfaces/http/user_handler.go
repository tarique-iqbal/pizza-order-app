package http

import (
	"api-service/internal/application/user"
	"api-service/internal/domain/auth"
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
	reqCtx := ctx.Request.Context()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("uniqueEmail", h.useCases.CustomValidator.UniqueEmail)
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		errors := validation.ExtractValidationErrors(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors})
		return
	}

	response, err := h.useCases.CreateUser.Execute(reqCtx, input)
	if err != nil {
		status := h.getHTTPStatusCode(err)
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (h *UserHandler) getHTTPStatusCode(err error) int {
	switch err {
	case auth.ErrCodeInvalid,
		auth.ErrCodeNotIssued:
		return http.StatusBadRequest
	case auth.ErrCodeExpired:
		return http.StatusGone
	case auth.ErrCodeUsed:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
