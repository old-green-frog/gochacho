CREATE TABLE IF NOT EXISTS equipment (
    id SERIAL PRIMARY KEY NOT NULL,
    brand VARCHAR(100),
    model VARCHAR(10),
    date_of_issue DATE
)
