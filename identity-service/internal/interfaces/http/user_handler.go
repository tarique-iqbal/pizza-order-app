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
	regCustomer *user.RegisterCustomer
	regOwner    *user.RegisterOwner
	findByID    *user.FindByID
}

func NewUserHandler(
	regCustomer *user.RegisterCustomer,
	regOwner *user.RegisterOwner,
	findByID *user.FindByID,
) *UserHandler {
	return &UserHandler{
		regCustomer: regCustomer,
		regOwner:    regOwner,
		findByID:    findByID,
	}
}

func (h *UserHandler) RegisterOwner(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	var input user.RegisterOwnerRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		errors := validation.ExtractValidationErrors(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors})
		return
	}

	response, err := h.regOwner.Execute(reqCtx, input)
	if err != nil {
		status := mapper.MapErrorToHTTPStatus(err)
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (h *UserHandler) RegisterCustomer(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	var input user.RegisterCustomerRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		errors := validation.ExtractValidationErrors(err)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errors})
		return
	}

	response, err := h.regCustomer.Execute(reqCtx, input)
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

	response, err := h.findByID.Execute(ctx, userID)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
