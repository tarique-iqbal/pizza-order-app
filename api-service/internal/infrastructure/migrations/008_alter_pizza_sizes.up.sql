-- Create composite unique index
CREATE UNIQUE INDEX idx_unique_restaurant_size
ON pizza_sizes (restaurant_id, size);

-- Drop the old one if it exists
DROP INDEX IF EXISTS idx_pizza_sizes_restaurant_id;
