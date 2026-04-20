-- Drop the index on restaurant_id (if it exists)
DROP INDEX IF EXISTS idx_pizza_sizes_restaurant_id;

-- Drop the pizza_sizes table (will also drop its constraints)
DROP TABLE IF EXISTS pizza_sizes;
