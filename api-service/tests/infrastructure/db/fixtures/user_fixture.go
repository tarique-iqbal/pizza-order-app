package fixtures

import (
	"api-service/internal/domain/user"
	"api-service/internal/infrastructure/security"
	"time"

	"gorm.io/gorm"
)

func LoadUserFixtures(db *gorm.DB) error {
	hasher := security.NewPasswordHasher()
	password, _ := hasher.Hash("plainPassword")

	users := []user.User{
		{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Password:  password,
			Role:      "User",
			CreatedAt: time.Now(),
		},
	}

	for _, u := range users {
		db.Create(&u)
	}

	return nil
}

func CreateUser(db *gorm.DB, role string) (*user.User, error) {
	user := user.User{
		FirstName: "Sofia",
		LastName:  "Harland",
		Email:     "sofia.harland@example.com",
		Password:  "hashedpassword",
		Role:      role,
		CreatedAt: time.Now(),
	}
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
