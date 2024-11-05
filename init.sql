CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(100) NOT NULL,
    token VARCHAR(100) DEFAULT ''
);

CREATE TABLE IF NOT EXISTS documents (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(255) NOT NULL,
    content BYTEA,
    public BOOLEAN DEFAULT FALSE,
    "grant" TEXT[],
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);