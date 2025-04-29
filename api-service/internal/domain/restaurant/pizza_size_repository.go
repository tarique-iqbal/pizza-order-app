package restaurant

import (
	"context"
)

type PizzaSizeRepository interface {
	Create(ctx context.Context, size *PizzaSize) error
}
