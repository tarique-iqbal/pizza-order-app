package restaurant

import "context"

type Geocoder interface {
	GeocodeAddress(ctx context.Context, address RestaurantAddress) (lat float64, lng float64, err error)
}
