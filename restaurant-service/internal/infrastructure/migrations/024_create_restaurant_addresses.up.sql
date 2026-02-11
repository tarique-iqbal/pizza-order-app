CREATE TABLE restaurant_addresses (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    house VARCHAR(63) NOT NULL,
    street VARCHAR(127) NOT NULL,
    city VARCHAR(63) NOT NULL,
    state VARCHAR(63),
    postal_code VARCHAR(10) NOT NULL,
    full_text TEXT NOT NULL,
    lat DOUBLE PRECISION,
    lon DOUBLE PRECISION
);

CREATE INDEX idx_restaurant_addresses_fulltext
ON restaurant_addresses
USING gin(to_tsvector('german', full_text));
