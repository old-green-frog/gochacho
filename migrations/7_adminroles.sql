CREATE TABLE IF NOT EXISTS adminroles (
    id SERIAL PRIMARY KEY NOT NULL,
    val VARCHAR(50)
);

INSERT INTO adminroles(id, val) VALUES 
(1, 'Руководитель'),
(2, 'Механик'),
(3, 'Кассир-контролер')
ON CONFLICT DO NOTHING;
