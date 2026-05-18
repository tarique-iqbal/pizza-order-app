-- global shared table
CREATE TABLE pizza_sizes (
    id UUID PRIMARY KEY,
    diameter_cm SMALLINT NOT NULL
        CONSTRAINT ck_pizza_sizes_diameter_cm
        CHECK (diameter_cm BETWEEN 20 AND 45),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT uq_pizza_sizes_diameter_cm
        UNIQUE (diameter_cm)
);

INSERT INTO pizza_sizes (
    id,
    diameter_cm
)
VALUES
('019e3c11-b1db-714f-bdff-c003d58b7377', 20),
('019e3c11-b1db-7a03-a3d1-5dcdc69daa24', 22),
('019e3c11-b1db-75f7-8057-419bf0ed7181', 24),
('019e3c11-b1db-78ed-b8e9-6e3e1672072f', 26),
('019e3c11-b1db-7335-b63e-93b529beac8a', 28),
('019e3c11-b1db-7793-896c-6c2c560a282f', 30),
('019e3c11-b1db-7c58-bdd7-42caf426d47d', 32),
('019e3c11-b1db-7016-8af2-b17d90fc08f6', 34),
('019e3c11-b1db-72fa-9c7f-bce2bad95655', 36),
('019e3c11-b1db-7114-82dd-ce6b814a5b98', 38),
('019e3c11-b1db-718d-b679-3cf231c03942', 40),
('019e3c11-b1db-7364-9ea7-4aa469e30dca', 42),
('019e3c11-b1db-7713-b56d-67eb1eeadcf3', 45)
ON CONFLICT (diameter_cm) DO NOTHING;
