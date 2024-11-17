CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password_salt BYTEA NOT NULL,
    password_hash BYTEA NOT NULL
);

-- Create a unique index on the username column to enforce uniqueness
CREATE UNIQUE INDEX users_username_idx ON users(LOWER(username));

CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(255) NOT NULL,
    uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    content BYTEA NOT NULL,
    size INTEGER NOT NULL,  -- Consider using BIGINT if you expect large files
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX files_uploaded_at_idx ON files(uploaded_at);
CREATE INDEX files_size_idx ON files(size);
CREATE INDEX files_name_idx ON files(LOWER(name));