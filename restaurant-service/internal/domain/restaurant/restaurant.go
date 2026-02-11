package restaurant

import (
	"time"

	"github.com/google/uuid"
)

type Restaurant struct {
	ID             uint              `gorm:"primaryKey"`
	RestaurantUUID uuid.UUID         `gorm:"type:uuid;not null;uniqueIndex"`
	UserID         uint              `gorm:"not null;index"`
	Name           string            `gorm:"size:200;not null"`
	Slug           string            `gorm:"size:255;uniqueIndex;not null"`
	Email          string            `gorm:"size:255;uniqueIndex;not null"`
	Phone          string            `gorm:"size:32;not null"`
	AddressID      uint              `gorm:"not null"`
	Address        RestaurantAddress `gorm:"foreignKey:AddressID"`
	DeliveryType   string            `gorm:"size:32;not null"`
	DeliveryKm     int               `gorm:"not null;"`
	Specialties    string            `gorm:"type:text"`
	CreatedAt      time.Time         `gorm:"autoCreateTime"`
	UpdatedAt      *time.Time        `gorm:"autoUpdateTime;default:null"`
}

func (Restaurant) TableName() string {
	return "restaurants"
}
