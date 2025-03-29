package restaurant

type RestaurantCreateDTO struct {
	UserID  uint   `json:"user_id"`
	Name    string `json:"name" binding:"required"`
	Slug    string `json:"slug" binding:"required,uniqueRSlug"`
	Address string `json:"address" binding:"required"`
}

type RestaurantResponseDTO struct {
	ID        uint    `json:"id"`
	UserID    uint    `json:"user_id"`
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	Address   string  `json:"address"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}
