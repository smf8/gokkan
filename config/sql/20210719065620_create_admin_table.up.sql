CREATE TABLE IF NOT EXISTS admins(
    id bigserial PRIMARY KEY,
    username VARCHAR (50) UNIQUE NOT NULL,
    password VARCHAR (50) NOT NULL
);