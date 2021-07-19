CREATE TABLE IF NOT EXISTS items(
    id BIGSERIAL PRIMARY KEY,
    category_id BIGSERIAL DEFAULT 1,
    price FLOAT(8),
    remaining INTEGER,
    sold INTEGER,
    photo_url varchar(255),
    created_at      timestamp not null default now(),
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE SET DEFAULT
 );