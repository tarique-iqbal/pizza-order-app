package db

import (
	"api-service/internal/domain/restaurant"
	"api-service/internal/domain/user"
	"log"

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

	sqlDB, sqlErr := db.DB()
	if sqlErr != nil {
		log.Fatalf("Failed to get *sql.DB from GORM: %v", sqlErr)
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	db.Exec("PRAGMA foreign_keys = ON")

	err := db.AutoMigrate(
		&user.User{},
		&restaurant.Restaurant{},
		&restaurant.PizzaSize{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}
