-- +goose Up
CREATE TABLE IF NOT EXISTS issues (
    id SERIAL PRIMARY KEY,
    error TEXT NOT NULL,
    cause TEXT NOT NULL,
    fix TEXT NOT NULL
);

-- +goose Down
DROP TABLE issues;