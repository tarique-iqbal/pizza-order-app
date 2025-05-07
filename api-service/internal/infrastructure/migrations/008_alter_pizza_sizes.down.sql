-- Drop the composite unique index
DROP INDEX IF EXISTS idx_unique_restaurant_size;

-- Recreate the old index on restaurant_id
CREATE INDEX idx_pizza_sizes_restaurant_id
ON pizza_sizes (restaurant_id);
