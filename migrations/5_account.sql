CREATE TABLE IF NOT EXISTS account (
    id SERIAL PRIMARY KEY NOT NULL,
    meter_id INT REFERENCES meter(id),
    customer_name VARCHAR(200) NOT NULL,
    customer_address VARCHAR,
    customer_number VARCHAR
)