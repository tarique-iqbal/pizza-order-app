ALTER TABLE restaurant_addresses ADD COLUMN full_text VARCHAR(511) NOT NULL;

CREATE INDEX idx_restaurant_addresses_full_text 
ON restaurant_addresses 
USING GIN (to_tsvector('german', full_text));
