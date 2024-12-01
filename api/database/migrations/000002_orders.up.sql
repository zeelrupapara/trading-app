-- +migrate Up
CREATE TABLE orders (
    id VARCHAR (50) PRIMARY KEY,
    user_id VARCHAR (50) NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    volume DECIMAL(18, 8) NOT NULL,
    order_type VARCHAR(4) CHECK (order_type IN ('buy', 'sell')) NOT NULL,
    price VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
);
