CREATE TABLE IF NOT EXISTS categories(
         id bigserial PRIMARY KEY,
         name VARCHAR(255) NOT NULL UNIQUE
    );

INSERT INTO categories VALUES (1, 'default');