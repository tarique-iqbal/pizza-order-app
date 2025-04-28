package restaurant

import (
	"time"
)

type PizzaSize struct {
	ID           uint       `gorm:"primaryKey"`
	RestaurantID uint       `gorm:"not null;index"`
	Title        string     `gorm:"size:63;not null"`
	Size         int        `gorm:"not null"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	UpdatedAt    *time.Time `gorm:"autoUpdateTime;default:null"`
}

func (PizzaSize) TableName() string {
	return "pizza_sizes"
}
