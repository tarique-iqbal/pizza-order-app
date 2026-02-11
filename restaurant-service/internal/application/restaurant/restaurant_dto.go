package restaurant

type RestaurantCreateDTO struct {
	UserID       uint     `json:"user_id"`
	Name         string   `json:"name" binding:"required,max=100"`
	Email        string   `json:"email" binding:"required,email"`
	Phone        string   `json:"phone" binding:"required"`
	House        string   `json:"house" binding:"required,max=63"`
	Street       string   `json:"street" binding:"required,max=127"`
	City         string   `json:"city" binding:"required,alphaunicode,max=63"`
	PostalCode   string   `json:"postal_code" binding:"required"`
	DeliveryType string   `json:"delivery_type" binding:"required,oneof=pick_up own_delivery third_party"`
	DeliveryKm   int      `json:"delivery_km" binding:"required,min=1,max=25"`
	Specialties  []string `json:"specialties"`
}

type RestaurantResponseDTO struct {
	ID           uint     `json:"id"`
	UserID       uint     `json:"user_id"`
	Name         string   `json:"name"`
	Slug         string   `json:"slug"`
	Email        string   `json:"email"`
	Phone        string   `json:"phone"`
	Address      string   `json:"address"`
	DeliveryType string   `json:"delivery_type"`
	DeliveryKm   int      `json:"delivery_km"`
	Specialties  []string `json:"specialties"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    *string  `json:"updated_at,omitempty"`
}
