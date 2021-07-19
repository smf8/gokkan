CREATE TABLE IF NOT EXISTS items(
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT DEFAULT 1,
    price FLOAT(8) NOT NULL,
    remaining INTEGER NOT NULL,
    sold INTEGER DEFAULT 0,
    photo_url varchar(255),
    created_at      timestamp NOT NULL DEFAULT now(),
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE SET DEFAULT
 );