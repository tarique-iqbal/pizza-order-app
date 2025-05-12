-- Create the enum type only if it doesn't already exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'restaurant_delivery_type_enum'
    ) THEN
        CREATE TYPE restaurant_delivery_type_enum AS ENUM ('pick_up', 'own_delivery', 'third_party');
    END IF;
END$$;

-- ALTER TABLE for changes
ALTER TABLE restaurants
    DROP COLUMN address,
    ADD COLUMN address_id INT REFERENCES restaurant_addresses(id),
    ADD COLUMN phone VARCHAR(32) NOT NULL,
    ADD COLUMN email VARCHAR(255) NOT NULL UNIQUE,
    ADD COLUMN delivery_type restaurant_delivery_type_enum NOT NULL DEFAULT 'pick_up',
    ADD COLUMN specialties VARCHAR(255);
