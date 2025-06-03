-- Create the enum type only if it doesn't already exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'restaurant_delivery_type_enum'
    ) THEN
        CREATE TYPE restaurant_delivery_type_enum AS ENUM ('pick_up', 'own_delivery', 'third_party');
    END IF;
END$$;

-- Create restaurants table
CREATE TABLE restaurants (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    restaurant_uuid UUID NOT NULL UNIQUE,
    user_id INTEGER NOT NULL,
    address_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(32) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    delivery_type restaurant_delivery_type_enum NOT NULL DEFAULT 'pick_up',
    delivery_km SMALLINT NOT NULL CONSTRAINT ck_delivery_km CHECK (delivery_km BETWEEN 1 AND 25),
    specialties VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    CONSTRAINT fk_restaurants_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_restaurants_address FOREIGN KEY (address_id) REFERENCES restaurant_addresses(id)
);
