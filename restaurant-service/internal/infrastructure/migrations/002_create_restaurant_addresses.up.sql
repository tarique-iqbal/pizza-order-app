CREATE TABLE restaurant_addresses (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    restaurant_id UUID NOT NULL UNIQUE,
    house VARCHAR(64) NOT NULL,
    street VARCHAR(128) NOT NULL,
    city VARCHAR(64) NOT NULL,
    state VARCHAR(64),
    postal_code VARCHAR(10) NOT NULL,
    full_text TEXT NOT NULL,
    lat DOUBLE PRECISION,
    lon DOUBLE PRECISION,

    CONSTRAINT fk_restaurant_addresses_restaurant
        FOREIGN KEY (restaurant_id)
        REFERENCES restaurants(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_restaurant_addresses_fulltext
ON restaurant_addresses
USING gin (to_tsvector('german', full_text));

CREATE INDEX idx_restaurant_addresses_restaurant_id
ON restaurant_addresses (restaurant_id);
