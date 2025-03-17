package user_test

import (
	applicationUser "pizza-order-api/internal/application/user"
	"pizza-order-api/internal/domain/user"
	"pizza-order-api/internal/infrastructure/persistence"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&user.User{})
	return db
}

func TestCreateUserUseCase(t *testing.T) {
	db := setupTestDB()
	userRepo := persistence.NewUserRepository(db)
	useCase := applicationUser.NewCreateUserUseCase(userRepo)

	input := applicationUser.UserCreateDTO{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
		Password:  "securepassword",
	}

	user, err := useCase.Execute(input)

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Jane", user.FirstName)
}
