-- Revert role column back to VARCHAR(20)
ALTER TABLE users
  ALTER COLUMN role DROP DEFAULT,
  ALTER COLUMN role TYPE VARCHAR(20) USING LOWER(role),
  ALTER COLUMN role SET DEFAULT 'user';

-- Drop ENUM type (only if no other columns are using it)
DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role_enum') THEN
        DROP TYPE user_role_enum;
    END IF;
END$$;
