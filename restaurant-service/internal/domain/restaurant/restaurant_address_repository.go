package restaurant

import (
	"context"

	"gorm.io/gorm"
)

type RestaurantAddressRepository interface {
	WithTx(tx *gorm.DB) RestaurantAddressRepository
	Create(ctx context.Context, address *RestaurantAddress) error
}
