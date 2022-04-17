CREATE TABLE IF NOT EXISTS admins (
    serial_num SERIAL PRIMARY KEY NOT NULL,
    role_id INT REFERENCES adminroles(id)
)
