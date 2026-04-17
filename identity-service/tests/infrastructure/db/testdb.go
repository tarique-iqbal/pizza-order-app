package db

import (
	"identity-service/internal/infrastructure/db"
	"log"

	"gorm.io/gorm"
)

func SetupTestDB() *gorm.DB {
	tdb, err := db.InitDB()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	tdb.Exec(`
		DROP TABLE IF EXISTS email_verifications CASCADE;
		DROP TABLE IF EXISTS users CASCADE;

		DROP TYPE IF EXISTS user_status_enum CASCADE;
		DROP TYPE IF EXISTS user_role_enum CASCADE;

		CREATE TYPE user_status_enum AS ENUM ('active', 'inactive', 'suspended');
		CREATE TYPE user_role_enum AS ENUM ('user', 'owner', 'admin');

		CREATE TABLE users (
			id UUID PRIMARY KEY,
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			role user_role_enum DEFAULT 'user'::user_role_enum,
			status user_status_enum DEFAULT 'active'::user_status_enum,
			logged_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP
		);

		CREATE TABLE email_verifications (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) NOT NULL,
			code CHAR(6) NOT NULL,
			is_used BOOLEAN DEFAULT FALSE,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)

	return tdb
}
