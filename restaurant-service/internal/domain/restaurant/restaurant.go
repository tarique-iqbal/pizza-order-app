package restaurant

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type DeliveryType string
type RestaurantStatus string

const (
	DeliveryOwn      DeliveryType = "own"
	DeliveryExternal DeliveryType = "external"
	DeliveryNone     DeliveryType = "none"
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
	Specialties  datatypes.JSON   `gorm:"type:jsonb;default:'[]';not null"`
	Checklist    Checklist        `gorm:"type:jsonb;serializer:json;not null"`
	Status       RestaurantStatus `gorm:"type:restaurant_status_enum;not null;default:'draft'"`
	CreatedAt    time.Time        `gorm:"type:timestamptz;not null"`
	UpdatedAt    *time.Time       `gorm:"type:timestamptz"`
}

func (Restaurant) TableName() string {
	return "restaurants"
}
