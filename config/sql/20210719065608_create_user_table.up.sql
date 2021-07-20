CREATE TABLE IF NOT EXISTS users(
    id bigserial PRIMARY KEY,
    username VARCHAR (50) UNIQUE NOT NULL,
    password VARCHAR (255) NOT NULL,
    full_name VARCHAR(255),
    billing_address VARCHAR (300),
    balance FLOAT(8) DEFAULT 0
    );

CREATE INDEX username_idx ON users (username);