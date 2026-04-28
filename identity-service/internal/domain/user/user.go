package user

import (
	"time"

	"github.com/google/uuid"
)

const DefaultStatus = "active"

type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	FirstName string     `gorm:"size:255;not null"`
	LastName  string     `gorm:"size:255;not null"`
	Email     string     `gorm:"size:255;unique;not null"`
	Password  string     `gorm:"not null"`
	Role      string     `gorm:"type:user_role_enum;default:'user'"`
	Status    string     `gorm:"type:user_status_enum;default:'active'"`
	LoggedAt  *time.Time `gorm:"column:logged_at;default:null"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;default:null"`
}

func (User) TableName() string {
	return "users"
}
