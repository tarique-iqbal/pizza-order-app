package fixtures

import (
	"pizza-order-api/internal/domain/user"
	"time"

	"gorm.io/gorm"
)

func LoadUserFixtures(db *gorm.DB) error {
	users := []user.User{
		{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Password:  "hashedpassword",
			Role:      "user",
			CreatedAt: time.Now(),
			UpdatedAt: nil,
		},
	}

	for _, u := range users {
		db.Create(&u)
	}

	return nil
}
