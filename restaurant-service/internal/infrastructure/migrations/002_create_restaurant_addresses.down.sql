-- Drop indexes first (explicit and symmetrical)
DROP INDEX IF EXISTS idx_restaurant_addresses_fulltext;
DROP INDEX IF EXISTS idx_restaurant_addresses_restaurant_id;

-- Drop the table (this also removes the FK constraint)
DROP TABLE IF EXISTS restaurant_addresses;
