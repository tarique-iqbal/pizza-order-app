package restaurant

import (
	"context"

	"gorm.io/gorm"
)

type RestaurantAddressRepository interface {
	WithTx(tx *gorm.DB) RestaurantAddressRepository
	Create(ctx context.Context, address *RestaurantAddress) error
}

type RestaurantRepository interface {
	WithTx(tx *gorm.DB) RestaurantRepository
	Create(ctx context.Context, r *Restaurant) error
	FindBySlug(ctx context.Context, slug string) (*Restaurant, error)
	IsOwnedBy(ctx context.Context, restaurantID uint, ownerID uint) (bool, error)
	IsSlugExists(ctx context.Context, slug string) (bool, error)
	IsEmailExists(ctx context.Context, email string) (bool, error)
}

type PizzaSizeRepository interface {
	Create(ctx context.Context, size *PizzaSize) error
	PizzaSizeExists(ctx context.Context, restaurantID uint, size int) (bool, error)
}
