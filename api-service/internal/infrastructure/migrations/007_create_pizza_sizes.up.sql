CREATE TABLE pizza_sizes (
    id SERIAL PRIMARY KEY,
    restaurant_id INT NOT NULL,
    title VARCHAR(63) NOT NULL,
    size INTEGER NOT NULL CHECK (size > 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT NULL,
    CONSTRAINT fk_pizza_sizes_restaurants FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE
);

CREATE INDEX idx_pizza_sizes_restaurant_id ON pizza_sizes (restaurant_id);
