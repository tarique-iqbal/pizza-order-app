package http

import (
	"identity-service/internal/application/user"
	"identity-service/internal/interfaces/http/mapper"
	"identity-service/internal/interfaces/http/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	register *user.Register
	findByID *user.FindByID
}

func NewUserHandler(reg *user.Register, findByID *user.FindByID) *UserHandler {
	return &UserHandler{register: reg, findByID: findByID}
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

func (h *UserHandler) FindByID(ctx *gin.Context) {
	idParam := ctx.Param("id")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	res, err := h.findByID.Execute(ctx, userID)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
