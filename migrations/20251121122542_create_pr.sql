-- +goose Up
-- +goose StatementBegin
CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE IF NOT EXISTS pull_requests (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    author_id  TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status     pr_status DEFAULT 'OPEN',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    merged_at  TIMESTAMP WITH TIME ZONE
);

-- Связь многие ко многим
CREATE TABLE IF NOT EXISTS pr_reviewers (
    pr_id       TEXT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_removed  BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (pr_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pull_requests;
DROP TABLE IF EXISTS pr_reviewers;
-- +goose StatementEnd
