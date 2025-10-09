CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY  NOT NULL DEFAULT gen_random_uuid(),
   first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255),
    password VARCHAR(255),
    email VARCHAR(255),
   last_login TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL
);