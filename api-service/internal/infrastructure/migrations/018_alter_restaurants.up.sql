ALTER TABLE restaurants
    ADD COLUMN delivery_km SMALLINT NOT NULL CHECK (delivery_km BETWEEN 1 AND 25);
