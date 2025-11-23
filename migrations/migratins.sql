DROP TABLE IF EXISTS pull_requests CASCADE;
DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE users (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    team_name   TEXT NOT NULL,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE
);


CREATE TABLE pull_requests (
    id            TEXT PRIMARY KEY,
    title         TEXT NOT NULL,
    description   TEXT NOT NULL DEFAULT '',
    author_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    team_id       TEXT NOT NULL,
    reviewers     TEXT[] NOT NULL DEFAULT '{}',
    status        TEXT NOT NULL CHECK (status IN ('open', 'merged')),
    merged_at     TEXT
);
