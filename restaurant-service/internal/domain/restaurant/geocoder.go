package restaurant

type GeocoderService interface {
	GeocodeAddress(address RestaurantAddress) (lat float64, lng float64, err error)
}
