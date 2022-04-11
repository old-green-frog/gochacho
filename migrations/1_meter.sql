CREATE TABLE IF NOT EXISTS meter (
    id SERIAL PRIMARY KEY NOT NULL,
    date_of_plug DATE,
    price_for_plug DECIMAL(10, 2),
    date_of_check DATE,
    price_for_check DECIMAL(10, 2)
)