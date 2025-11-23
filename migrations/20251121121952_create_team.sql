-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS teams (
    team_name   TEXT PRIMARY KEY,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

ALTER TABLE users 
ADD COLUMN team_name TEXT REFERENCES teams(team_name) 
ON DELETE SET NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN team_name;
DROP TABLE IF EXISTS teams;
-- +goose StatementEnd
