package auth

import (
	"time"
)

type EmailVerification struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"size:255;not null;index"`
	Code      string    `gorm:"type:char(6);not null"`
	IsUsed    bool      `gorm:"default:false"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (EmailVerification) TableName() string {
	return "email_verifications"
}
