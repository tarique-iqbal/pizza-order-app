package restaurant

import (
	"api-service/internal/domain/restaurant"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type CreateRestaurantUseCase struct {
	db           *gorm.DB
	repo         restaurant.RestaurantRepository
	restAddrRepo restaurant.RestaurantAddressRepository
}

func NewCreateRestaurantUseCase(
	db *gorm.DB,
	repo restaurant.RestaurantRepository,
	addressRepo restaurant.RestaurantAddressRepository,
) *CreateRestaurantUseCase {
	return &CreateRestaurantUseCase{
		db:           db,
		repo:         repo,
		restAddrRepo: addressRepo,
	}
}

func (uc *CreateRestaurantUseCase) Execute(ctx context.Context, input RestaurantCreateDTO) (RestaurantResponseDTO, error) {
	if err := uc.checkEmailUnique(ctx, input.Email); err != nil {
		return RestaurantResponseDTO{}, err
	}

	slugFinal, err := uc.generateUniqueSlug(ctx, input.Name, input.City)
	if err != nil {
		return RestaurantResponseDTO{}, err
	}

	address := uc.buildAddress(input)

	newRestaurant, newAddress, err := uc.saveAddressAndRestaurant(ctx, address, input, slugFinal)
	if err != nil {
		return RestaurantResponseDTO{}, err
	}

	return RestaurantResponseDTO{
		ID:           newRestaurant.ID,
		UserID:       newRestaurant.UserID,
		Name:         newRestaurant.Name,
		Slug:         newRestaurant.Slug,
		Email:        newRestaurant.Email,
		Phone:        newRestaurant.Phone,
		Address:      newAddress.FullText,
		DeliveryType: newRestaurant.DeliveryType,
		DeliveryKm:   newRestaurant.DeliveryKm,
		Specialties:  input.Specialties,
		CreatedAt:    newRestaurant.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (uc *CreateRestaurantUseCase) checkEmailUnique(ctx context.Context, email string) error {
	exists, err := uc.repo.IsEmailExists(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return restaurant.ErrEmailAlreadyExists
	}
	return nil
}

func (uc *CreateRestaurantUseCase) generateUniqueSlug(ctx context.Context, name, city string) (string, error) {
	slugBase := slug.Make(fmt.Sprintf("%s-%s", name, city))
	slugFinal := slugBase
	attempt := 1
	for {
		exists, err := uc.repo.IsSlugExists(ctx, slugFinal)
		if err != nil {
			return "", err
		}
		if !exists {
			break
		}
		attempt++
		slugFinal = fmt.Sprintf("%s-%d", slugBase, attempt)
	}
	return slugFinal, nil
}

func (uc *CreateRestaurantUseCase) buildAddress(input RestaurantCreateDTO) restaurant.RestaurantAddress {
	return restaurant.RestaurantAddress{
		House:      input.House,
		Street:     input.Street,
		City:       input.City,
		PostalCode: input.PostalCode,
		FullText: fmt.Sprintf("%s %s, %s %s",
			input.Street,
			input.House,
			input.PostalCode,
			input.City,
		),
	}
}

func (uc *CreateRestaurantUseCase) saveAddressAndRestaurant(
	ctx context.Context,
	address restaurant.RestaurantAddress,
	input RestaurantCreateDTO,
	slugFinal string,
) (restaurant.Restaurant, restaurant.RestaurantAddress, error) {
	var newRestaurant restaurant.Restaurant
	var newAddress restaurant.RestaurantAddress

	err := uc.db.Transaction(func(tx *gorm.DB) error {
		if err := uc.restAddrRepo.WithTx(tx).Create(ctx, &address); err != nil {
			return err
		}
		newAddress = address

		newRestaurant = restaurant.Restaurant{
			RestaurantUUID: uuid.New(),
			UserID:         input.UserID,
			Name:           input.Name,
			Slug:           slugFinal,
			Email:          input.Email,
			Phone:          input.Phone,
			AddressID:      address.ID,
			DeliveryType:   input.DeliveryType,
			DeliveryKm:     input.DeliveryKm,
			Specialties:    strings.Join(input.Specialties, ","),
		}
		if err := uc.repo.WithTx(tx).Create(ctx, &newRestaurant); err != nil {
			return err
		}
		return nil
	})
	return newRestaurant, newAddress, err
}
