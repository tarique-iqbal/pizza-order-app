package restaurant

import (
	"context"
)

type PizzaSizeRepository interface {
	Create(ctx context.Context, size *PizzaSize) error
	PizzaSizeExists(ctx context.Context, restaurantID uint, size int) (bool, error)
}
