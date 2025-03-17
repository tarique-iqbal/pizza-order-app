package persistence_test

import (
	"pizza-order-api/internal/domain/user"
	"pizza-order-api/internal/infrastructure/persistence"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&user.User{})
	return db
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB()
	repo := persistence.NewUserRepository(db)

	newUser := &user.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "hashedpassword",
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	err := repo.Create(newUser)

	assert.Nil(t, err)
	assert.NotZero(t, newUser.ID)
}
