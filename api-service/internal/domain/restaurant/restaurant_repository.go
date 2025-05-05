package restaurant

import "context"

type RestaurantRepository interface {
	Create(ctx context.Context, r *Restaurant) error
	FindBySlug(ctx context.Context, slug string) (*Restaurant, error)
	IsOwnedBy(ctx context.Context, restaurantID uint, ownerID uint) (bool, error)
	SlugExists(slug string) (bool, error)
}
