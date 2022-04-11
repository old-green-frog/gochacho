CREATE TABLE IF NOT EXISTS defectivelist (
    id SERIAL PRIMARY KEY NOT NULL,
    equipment_id INT REFERENCES equipment(id),
    defect VARCHAR(100),
    worktype_id INT REFERENCES worktype(id),
    brigade_id INT REFERENCES brigade(id),
    date_of_start DATE,
    desired_time INT,
    date_of_end DATE,
    date_of_accept DATE
)
