CREATE TABLE IF NOT EXISTS categories(
     id bigserial PRIMARY KEY,
     name VARCHAR(255)
);

INSERT INTO categories VALUES (1, 'default');