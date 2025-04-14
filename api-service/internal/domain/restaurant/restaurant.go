package restaurant

import "time"

type Restaurant struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id" gorm:"not null"`
	Name      string     `json:"name" gorm:"not null"`
	Slug      string     `json:"slug" gorm:"unique;not null"`
	Address   string     `json:"address" gorm:"not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"default:null"`
}
