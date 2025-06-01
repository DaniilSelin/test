CREATE SCHEMA IF NOT EXISTS %[1]s;

CREATE TABLE IF NOT EXISTS %[1]s.quotesbook (
    id    SERIAL PRIMARY KEY,
    author VARCHAR(255) NOT NULL,
    quote       TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Для быстрого поиска цитат по автору
CREATE INDEX IF NOT EXISTS idx_quotesbook_author
  ON %[1]s.quotesbook (author);