package restaurant

type PizzaSizeCreateDTO struct {
	Title string `json:"title" binding:"required"`
	Size  int    `json:"size" binding:"required,gt=0"`
}

type PizzaSizeResponseDTO struct {
	ID           uint    `json:"id"`
	RestaurantID uint    `json:"restaurant_id"`
	Title        string  `json:"title"`
	Size         int     `json:"size"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    *string `json:"updated_at,omitempty"`
}
