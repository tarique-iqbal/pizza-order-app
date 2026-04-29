package fixtures

import (
	"fmt"
	"identity-service/internal/domain/user"
	"identity-service/internal/infrastructure/security"
	"time"

	"github.com/google/uuid"
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
			Role:      "user",
			CreatedAt: time.Now().UTC(),
		},
		{
			FirstName: "Existing",
			LastName:  "User",
			Email:     "existing@example.com",
			Password:  password,
			Role:      "owner",
			CreatedAt: time.Now().UTC(),
		},
	}

	for _, u := range users {
		userID, _ := uuid.NewV7()
		u.ID = userID

		db.Create(&u)
	}

	return nil
}

func NewUser() *user.User {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	return &user.User{
		ID:        id,
		FirstName: "Sofia",
		LastName:  "Harland",
		Email:     randomEmail(),
		Password:  "hashedpassword",
		Role:      "user",
		CreatedAt: time.Now().UTC(),
	}
}

func CreateUser(db *gorm.DB, u *user.User) (*user.User, error) {
	if err := db.Create(u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func randomEmail() string {
	return fmt.Sprintf("user_%d@example.com", time.Now().UnixNano())
}
