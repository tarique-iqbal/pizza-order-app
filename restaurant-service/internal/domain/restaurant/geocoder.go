package restaurant

import "context"

type Geocoder interface {
	GeocodeAddress(ctx context.Context, address Address) (lat float64, lng float64, err error)
}
