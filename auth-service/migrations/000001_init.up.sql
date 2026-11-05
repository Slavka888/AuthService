CREATE TABLE IF NOT EXISTS users (
    id integer PRIMARY KEY,
    email text NOT NULL UNIQUE,
    password_hash blob NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

CREATE TABLE IF NOT EXISTS apps (
    id integer PRIMARY KEY,
    name text NOT NULL UNIQUE,
    secret text NOT NULL
);