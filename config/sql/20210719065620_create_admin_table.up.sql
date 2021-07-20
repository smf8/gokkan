CREATE TABLE IF NOT EXISTS admins(
    id bigserial PRIMARY KEY,
    username VARCHAR (50) UNIQUE NOT NULL,
    password VARCHAR (255) NOT NULL
);

INSERT INTO admins VALUES(1, 'admin', '8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918');