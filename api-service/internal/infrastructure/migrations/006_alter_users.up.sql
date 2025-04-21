-- Create ENUM types (only if they do not exist)
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role_enum') THEN
        CREATE TYPE user_role_enum AS ENUM ('User', 'Owner', 'Admin');
    END IF;
END$$;

-- Alter users table column role
ALTER TABLE users
  ALTER COLUMN role DROP DEFAULT,
  ALTER COLUMN role TYPE user_role_enum USING
    (INITCAP(role)::user_role_enum),
  ALTER COLUMN role SET DEFAULT 'User';
