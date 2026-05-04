package restaurant

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type DeliveryType string
type RestaurantStatus string

const (
	DeliveryOwn        DeliveryType = "own"
	DeliveryThirdParty DeliveryType = "third_party"
	DeliveryNone       DeliveryType = "none"
)

const (
	StatusDraft    RestaurantStatus = "draft"
	StatusReview   RestaurantStatus = "review"
	StatusActive   RestaurantStatus = "active"
	StatusInactive RestaurantStatus = "inactive"
	StatusDisabled RestaurantStatus = "disabled"
	StatusRejected RestaurantStatus = "rejected"
)

type Restaurant struct {
	ID           uuid.UUID        `gorm:"type:uuid;primaryKey"`
	OwnerID      uuid.UUID        `gorm:"type:uuid;not null;index"`
	Name         string           `gorm:"size:128;not null"`
	VATNumber    string           `gorm:"size:16;not null"`
	Slug         *string          `gorm:"size:255;uniqueIndex"`
	Email        *string          `gorm:"size:255;uniqueIndex"`
	Phone        *string          `gorm:"size:32"`
	Pickup       bool             `gorm:"not null;default:true"`
	DeliveryKm   *int16           `gorm:"check:delivery_km BETWEEN 1 AND 25"`
	DeliveryType DeliveryType     `gorm:"type:restaurant_delivery_type_enum;not null;default:'none'"`
	Specialties  *string          `gorm:"size:255"`
	Checklist    datatypes.JSON   `gorm:"type:jsonb;not null"`
	Status       RestaurantStatus `gorm:"type:restaurant_status_enum;not null;default:'draft'"`
	CreatedAt    time.Time        `gorm:"autoCreateTime"`
	UpdatedAt    *time.Time       `gorm:"autoUpdateTime"`
}

func (Restaurant) TableName() string {
	return "restaurants"
}
