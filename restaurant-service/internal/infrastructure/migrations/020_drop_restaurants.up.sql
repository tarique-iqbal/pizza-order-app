-- Drop foreign keys from tables that reference restaurants
ALTER TABLE pizza_sizes DROP CONSTRAINT IF EXISTS fk_pizza_sizes_restaurants;

-- Drop the restaurants table
DROP TABLE IF EXISTS restaurants;

-- Drop the restaurant_addresses table
DROP TABLE IF EXISTS restaurant_addresses;

-- Drop the enum type (optional — only if no other table uses it)
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'restaurant_delivery_type_enum'
    ) THEN
        DROP TYPE restaurant_delivery_type_enum;
    END IF;
END$$;
