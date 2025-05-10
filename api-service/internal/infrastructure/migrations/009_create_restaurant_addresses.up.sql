CREATE TABLE restaurant_addresses (
    id SERIAL PRIMARY KEY,
    street VARCHAR(255) NOT NULL,
    city VARCHAR(127) NOT NULL,
    state VARCHAR(63),
    postal_code VARCHAR(10) NOT NULL,
    lat DOUBLE PRECISION,
    lon DOUBLE PRECISION
);
