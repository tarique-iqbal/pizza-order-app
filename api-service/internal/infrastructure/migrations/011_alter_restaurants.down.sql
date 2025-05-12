-- Reverse ALTER TABLE: remove added columns, restore dropped column
ALTER TABLE restaurants
    DROP COLUMN specialties,
    DROP COLUMN delivery_type,
    DROP COLUMN email,
    DROP COLUMN phone,
    DROP COLUMN address_id,
    ADD COLUMN address VARCHAR(511) NOT NULL;

-- Drop enum type if exists and not in use
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'restaurant_delivery_type_enum'
    ) THEN
        -- Ensure no other table is using this type before dropping
        DROP TYPE restaurant_delivery_type_enum;
    END IF;
END$$;
