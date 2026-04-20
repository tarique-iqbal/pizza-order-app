-- Drop the restaurants table first to remove dependency
DROP TABLE IF EXISTS restaurants;

-- Then drop the GIN index on restaurant_addresses
DROP INDEX IF EXISTS idx_restaurant_addresses_fulltext;

-- Finally, drop restaurant_addresses
DROP TABLE IF EXISTS restaurant_addresses;
