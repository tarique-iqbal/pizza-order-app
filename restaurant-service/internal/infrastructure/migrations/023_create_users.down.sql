-- First, drop the table (removes columns and constraints)
DROP TABLE IF EXISTS users;

-- Then, drop the ENUM types if they exist
DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role_enum') THEN
        DROP TYPE user_role_enum;
    END IF;
END$$;

DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_status_enum') THEN
        DROP TYPE user_status_enum;
    END IF;
END$$;
