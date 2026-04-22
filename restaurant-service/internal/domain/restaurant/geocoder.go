package restaurant

type Geocoder interface {
	GeocodeAddress(address RestaurantAddress) (lat float64, lng float64, err error)
}
