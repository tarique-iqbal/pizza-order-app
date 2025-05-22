package restaurant

type RestaurantAddress struct {
	ID         uint    `gorm:"primaryKey"`
	House      string  `gorm:"size:63;not null"`
	Street     string  `gorm:"size:127;not null"`
	City       string  `gorm:"size:63;not null"`
	PostalCode string  `gorm:"size:10;not null"`
	FullText   string  `gorm:"type:text;not null"`
	Lat        float64 `gorm:"type:double precision"`
	Lon        float64 `gorm:"type:double precision"`
}

func (RestaurantAddress) TableName() string {
	return "restaurant_addresses"
}
