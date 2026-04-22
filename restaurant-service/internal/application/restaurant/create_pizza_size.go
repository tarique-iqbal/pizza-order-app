package restaurant

import (
	"context"
	"fmt"
	"time"

	"restaurant-service/internal/domain/restaurant"
	apperr "restaurant-service/internal/shared/errors"
)

type CreatePizzaSize struct {
	pizzaSizeRepo  restaurant.PizzaSizeRepository
	restaurantRepo restaurant.RestaurantRepository
}

func NewCreatePizzaSize(
	pizzaSizeRepo restaurant.PizzaSizeRepository,
	restaurantRepo restaurant.RestaurantRepository,
) *CreatePizzaSize {
	return &CreatePizzaSize{pizzaSizeRepo: pizzaSizeRepo, restaurantRepo: restaurantRepo}
}

func (uc *CreatePizzaSize) Execute(
	ctx context.Context,
	restaurantID uint,
	ownerID uint,
	input CreatePizzaSizeRequest,
) (PizzaSizeResponse, error) {
	owns, err := uc.restaurantRepo.IsOwnedBy(ctx, restaurantID, ownerID)
	if err != nil {
		return PizzaSizeResponse{}, fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !owns {
		return PizzaSizeResponse{}, apperr.ErrForbidden
	}

	exists, err := uc.pizzaSizeRepo.PizzaSizeExists(ctx, restaurantID, input.Size)
	if err != nil {
		return PizzaSizeResponse{}, err
	}
	if exists {
		return PizzaSizeResponse{}, restaurant.ErrPizzaSizeAlreadyExists
	}

	newPizzaSize := &restaurant.PizzaSize{
		RestaurantID: restaurantID,
		Title:        input.Title,
		Size:         input.Size,
	}

	if err := uc.pizzaSizeRepo.Create(ctx, newPizzaSize); err != nil {
		return PizzaSizeResponse{}, err
	}

	response := PizzaSizeResponse{
		ID:           newPizzaSize.ID,
		RestaurantID: newPizzaSize.RestaurantID,
		Title:        newPizzaSize.Title,
		Size:         newPizzaSize.Size,
		CreatedAt:    newPizzaSize.CreatedAt.Format(time.RFC3339),
	}

	return response, nil
}
