-- Create the enum type only if it doesn't already exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'restaurant_delivery_type_enum'
    ) THEN
        CREATE TYPE restaurant_delivery_type_enum AS ENUM (
            'own', 'third_party', 'none'
        );
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'restaurant_status_enum'
    ) THEN
        CREATE TYPE restaurant_status_enum AS ENUM (
            'draft', 'review', 'active', 'inactive', 'disabled', 'rejected'
        );
    END IF;
END$$;

-- Create restaurants table
CREATE TABLE restaurants (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL,
    name VARCHAR(128) NOT NULL,
    vat_number VARCHAR(16) NOT NULL,
    slug VARCHAR(255),
	email VARCHAR(255),
    phone VARCHAR(32),
    pickup BOOLEAN NOT NULL DEFAULT true,
    delivery_km SMALLINT
        CONSTRAINT ck_restaurants_delivery_km CHECK (delivery_km BETWEEN 1 AND 25),
	delivery_type restaurant_delivery_type_enum NOT NULL DEFAULT 'none',
    specialties VARCHAR(255),
    checklist JSONB NOT NULL,
    status restaurant_status_enum NOT NULL DEFAULT 'draft',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX ux_restaurants_slug
ON restaurants (slug)
WHERE slug IS NOT NULL;

CREATE UNIQUE INDEX ux_restaurants_email
ON restaurants (email)
WHERE email IS NOT NULL;

CREATE INDEX idx_restaurants_owner_id ON restaurants(owner_id);
