package restaurant

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RestaurantAddressRepository interface {
	WithTx(tx *gorm.DB) RestaurantAddressRepository
	Create(ctx context.Context, address *RestaurantAddress) error
}

type RestaurantRepository interface {
	WithTx(tx *gorm.DB) RestaurantRepository
	Create(ctx context.Context, res *Restaurant) error
	Update(ctx context.Context, res *Restaurant) error
	FindBySlug(ctx context.Context, slug string) (*Restaurant, error)
	FindByIDAndOwner(ctx context.Context, restaurantID uuid.UUID, ownerID uuid.UUID) (*Restaurant, error)
	IsOwner(ctx context.Context, restaurantID uuid.UUID, ownerID uuid.UUID) (bool, error)
	IsSlugExists(ctx context.Context, slug string) (bool, error)
	IsEmailExists(ctx context.Context, email string) (bool, error)
}

type PizzaSizeRepository interface {
	Create(ctx context.Context, size *PizzaSize) error
	PizzaSizeExists(ctx context.Context, restaurantID uuid.UUID, size int) (bool, error)
}
