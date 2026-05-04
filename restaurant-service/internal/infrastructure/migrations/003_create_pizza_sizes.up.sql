CREATE TABLE pizza_sizes (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    restaurant_id UUID NOT NULL,
    title VARCHAR(64) NOT NULL,
    size INTEGER NOT NULL CHECK (size > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,

    CONSTRAINT fk_pizza_sizes_restaurants
        FOREIGN KEY (restaurant_id)
        REFERENCES restaurants(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_pizza_sizes_restaurant_id_size
        UNIQUE (restaurant_id, size)
);
