package restaurant

import "context"

type RestaurantRepository interface {
	Create(r *Restaurant) error
	FindBySlug(slug string) (*Restaurant, error)
	IsOwnedBy(ctx context.Context, restaurantID uint, ownerID uint) (bool, error)
	SlugExists(slug string) (bool, error)
}
