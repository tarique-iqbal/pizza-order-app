-- Create ENUM types (only if they do not exist)
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_status_enum') THEN
        CREATE TYPE user_status_enum AS ENUM ('Active', 'Inactive', 'Suspended');
    END IF;
END$$;

DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'email_verified_enum') THEN
        CREATE TYPE email_verified_enum AS ENUM ('Yes', 'No');
    END IF;
END$$;

-- Alter users table to add new columns
ALTER TABLE users
ADD COLUMN IF NOT EXISTS status user_status_enum DEFAULT 'Active',
ADD COLUMN IF NOT EXISTS verified email_verified_enum DEFAULT 'No',
ADD COLUMN IF NOT EXISTS logged_at TIMESTAMP DEFAULT NULL;
