package restaurant

import "time"

type Restaurant struct {
	ID        uint       `gorm:"primaryKey"`
	UserID    uint       `gorm:"not null"`
	Name      string     `gorm:"size:255;not null"`
	Slug      string     `gorm:"size:255;unique;not null"`
	Address   string     `gorm:"size:511;not null"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;default:null"`
}
