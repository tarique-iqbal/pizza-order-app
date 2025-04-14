package restaurant

import (
	"api-service/internal/domain/restaurant"
	"time"
)

type CreateRestaurantUseCase struct {
	repo restaurant.RestaurantRepository
}

func NewCreateRestaurantUseCase(repo restaurant.RestaurantRepository) *CreateRestaurantUseCase {
	return &CreateRestaurantUseCase{repo: repo}
}

func (uc *CreateRestaurantUseCase) Execute(input RestaurantCreateDTO) (RestaurantResponseDTO, error) {
	newRestaurant := restaurant.Restaurant{
		UserID:    input.UserID,
		Name:      input.Name,
		Slug:      input.Slug,
		Address:   input.Address,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	err := uc.repo.Create(&newRestaurant)
	if err != nil {
		return RestaurantResponseDTO{}, err
	}

	response := RestaurantResponseDTO{
		ID:        newRestaurant.ID,
		UserID:    newRestaurant.UserID,
		Name:      newRestaurant.Name,
		Slug:      newRestaurant.Slug,
		Address:   newRestaurant.Address,
		CreatedAt: newRestaurant.CreatedAt.Format(time.RFC3339),
		UpdatedAt: nil,
	}

	return response, nil
}
