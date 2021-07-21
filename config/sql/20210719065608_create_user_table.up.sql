CREATE TABLE IF NOT EXISTS users(
    id bigserial PRIMARY KEY,
    username VARCHAR (50) UNIQUE NOT NULL,
    password VARCHAR (255) NOT NULL,
    full_name VARCHAR(255),
    billing_address VARCHAR (300),
    balance FLOAT(8) DEFAULT 0,
    is_admin BOOLEAN DEFAULT FALSE
    );

CREATE INDEX username_idx ON users (username);
INSERT INTO users(id, username, password, is_admin) VALUES(0, 'admin@admin.admin', '1d048b74d768cc1dc97c8ae37b2e15c6df1038aefea68e0a546d45a81373f068', true);