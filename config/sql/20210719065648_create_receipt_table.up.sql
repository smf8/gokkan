CREATE TABLE IF NOT EXISTS receipts(
    id BIGSERIAL PRIMARY KEY,
    item_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    quantity INTEGER NOT NULL,
    buyer_name VARCHAR(255) NOT NULL,
    billing_address VARCHAR(255) NOT NULL,
    price FLOAT(8) NOT NULL,
    date TIMESTAMP NOT NULL  DEFAULT now(),
    tracking_code VARCHAR(255) NOT NULL,
    status INTEGER NOT NULL,
    CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE NO ACTION,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE NO ACTION
);