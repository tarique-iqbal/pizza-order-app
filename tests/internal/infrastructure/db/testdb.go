package db

import (
	"log"
	"pizza-order-api/internal/domain/restaurant"
	"pizza-order-api/internal/domain/user"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to initialize test database: %v", err)
	}

	return db
}

func SetupTestDB() *gorm.DB {
	db := InitTestDB()

	err := db.AutoMigrate(
		&user.User{},
		&restaurant.Restaurant{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}
