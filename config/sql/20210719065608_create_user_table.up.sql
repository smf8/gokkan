CREATE TABLE IF NOT EXISTS users(
    id bigserial PRIMARY KEY,
    username VARCHAR (50) UNIQUE NOT NULL,
    password VARCHAR (50) NOT NULL,
    fullname VARCHAR(255),
    address VARCHAR (300),
    balance FLOAT(8) DEFAULT 0
    );