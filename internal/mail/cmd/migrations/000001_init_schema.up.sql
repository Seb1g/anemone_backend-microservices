CREATE TABLE temp_addresses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    address VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

CREATE TABLE emails (
    id SERIAL PRIMARY KEY,
    address_id INTEGER NOT NULL REFERENCES temp_addresses (id) ON DELETE CASCADE,
    sender VARCHAR(255) NOT NULL,
    recipients TEXT [] NOT NULL,
    subject TEXT,
    body TEXT,
    raw_data BYTEA,
    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);