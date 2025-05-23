package restaurant

import (
	"context"

	"gorm.io/gorm"
)

type RestaurantRepository interface {
	WithTx(tx *gorm.DB) RestaurantRepository
	Create(ctx context.Context, r *Restaurant) error
	FindBySlug(ctx context.Context, slug string) (*Restaurant, error)
	IsOwnedBy(ctx context.Context, restaurantID uint, ownerID uint) (bool, error)
	IsSlugExists(ctx context.Context, slug string) (bool, error)
	IsEmailExists(ctx context.Context, email string) (bool, error)
}
