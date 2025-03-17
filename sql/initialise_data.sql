CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    submitted_count INTEGER DEFAULT 1,
    last_checked_at TIMESTAMP NULL,
    last_response_time FLOAT NULL,
    is_valid BOOLEAN DEFAULT FALSE
    );

CREATE TABLE IF NOT EXISTS url_checks (
    id SERIAL PRIMARY KEY,
    url_id INTEGER REFERENCES urls(id) ON DELETE CASCADE,
    checked_at TIMESTAMP DEFAULT NOW(),
    response_time FLOAT NULL,
    status_code INTEGER NULL,
    is_success BOOLEAN NOT NULL
    );

CREATE INDEX IF NOT EXISTS idx_urls_last_checked ON urls (last_checked_at);
CREATE INDEX IF NOT EXISTS idx_url_checks_url_id ON url_checks (url_id);