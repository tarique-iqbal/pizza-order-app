-- Create ENUM type only if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'email_verified_enum'
    ) THEN
        CREATE TYPE email_verified_enum AS ENUM ('Yes', 'No');
    END IF;
END
$$;

-- Add 'verified' column to 'users' table using the enum, defaulting to 'No'
ALTER TABLE users
ADD COLUMN IF NOT EXISTS verified email_verified_enum DEFAULT 'No';
