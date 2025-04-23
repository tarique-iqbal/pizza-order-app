package user

import (
	"time"
)

type User struct {
	ID        uint       `gorm:"primaryKey"`
	FirstName string     `gorm:"size:255;not null"`
	LastName  string     `gorm:"size:255;not null"`
	Email     string     `gorm:"size:255;unique;not null"`
	Password  string     `gorm:"not null"`
	Role      string     `gorm:"size:16;default:'user'"`
	Status    string     `gorm:"type:user_status_enum;default:'Active'"`
	LoggedAt  *time.Time `gorm:"column:logged_at;default:null"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;default:null"`
}
