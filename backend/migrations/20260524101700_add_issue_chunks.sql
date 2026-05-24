-- +goose Up

CREATE TABLE issue_chunks (
    id SERIAL PRIMARY KEY,

    issue_id INTEGER REFERENCES issues(id) ON DELETE CASCADE,

    chunk_index INTEGER NOT NULL,

    content TEXT NOT NULL,

    embedding vector(768),

    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down

DROP TABLE issue_chunks;