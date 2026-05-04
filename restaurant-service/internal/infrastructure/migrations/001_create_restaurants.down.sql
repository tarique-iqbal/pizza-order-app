-- Drop the restaurants table first (it depends on the enum types)
DROP TABLE IF EXISTS restaurants;

-- Drop enum types if they exist
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'restaurant_delivery_type_enum'
    ) THEN
        DROP TYPE restaurant_delivery_type_enum;
    END IF;
END$$;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'restaurant_status_enum'
    ) THEN
        DROP TYPE restaurant_status_enum;
    END IF;
END$$;
