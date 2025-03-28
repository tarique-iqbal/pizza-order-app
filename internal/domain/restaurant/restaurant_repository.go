package restaurant

type RestaurantRepository interface {
	Create(r *Restaurant) error
	FindBySlug(slug string) (*Restaurant, error)
}
