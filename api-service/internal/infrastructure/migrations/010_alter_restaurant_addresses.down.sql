-- Step 1: Drop the GIN index
DROP INDEX IF EXISTS idx_restaurant_addresses_full_text;

-- Step 2: Drop the column
ALTER TABLE restaurant_addresses DROP COLUMN IF EXISTS full_text;
