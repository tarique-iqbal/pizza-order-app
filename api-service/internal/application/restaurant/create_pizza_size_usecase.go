package restaurant

import (
	"context"
	"fmt"
	"time"

	"api-service/internal/domain/restaurant"
	sErrors "api-service/internal/shared/errors"
)

type CreatePizzaSizeUseCase struct {
	pizzaSizeRepo  restaurant.PizzaSizeRepository
	restaurantRepo restaurant.RestaurantRepository
}

func NewCreatePizzaSizeUseCase(
	pizzaSizeRepo restaurant.PizzaSizeRepository,
	restaurantRepo restaurant.RestaurantRepository,
) *CreatePizzaSizeUseCase {
	return &CreatePizzaSizeUseCase{pizzaSizeRepo: pizzaSizeRepo, restaurantRepo: restaurantRepo}
}

func (uc *CreatePizzaSizeUseCase) Execute(
	ctx context.Context,
	restaurantID uint,
	ownerID uint,
	input PizzaSizeCreateDTO,
) (PizzaSizeResponseDTO, error) {
	owns, err := uc.restaurantRepo.IsOwnedBy(ctx, restaurantID, ownerID)
	if err != nil {
		return PizzaSizeResponseDTO{}, fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !owns {
		return PizzaSizeResponseDTO{}, sErrors.ErrForbidden
	}

	newPizzaSize := &restaurant.PizzaSize{
		RestaurantID: restaurantID,
		Title:        input.Title,
		Size:         input.Size,
	}

	if err := uc.pizzaSizeRepo.Create(ctx, newPizzaSize); err != nil {
		return PizzaSizeResponseDTO{}, err
	}

	response := PizzaSizeResponseDTO{
		ID:        newPizzaSize.ID,
		Title:     newPizzaSize.Title,
		Size:      newPizzaSize.Size,
		CreatedAt: newPizzaSize.CreatedAt.Format(time.RFC3339),
	}

	return response, nil
}
