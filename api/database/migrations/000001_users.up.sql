-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR (50) PRIMARY KEY,
    first_name VARCHAR (50) NOT NULL,
    last_name VARCHAR (50) NOT NULL,
    email VARCHAR (50) UNIQUE NOT NULL,
    password VARCHAR (100),
    roles TEXT,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
