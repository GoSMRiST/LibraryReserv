CREATE TABLE IF NOT EXISTS reservations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    author TEXT NOT NULL,
    title TEXT NOT NULL,
    taken_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    return_at TIMESTAMPTZ
);