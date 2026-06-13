-- +goose Up
CREATE TABLE IF NOT EXISTS incidents (
    id SERIAL PRIMARY KEY,
    title TEXT,
    summary TEXT,
    symptoms TEXT[] DEFAULT '{}',
    evidence TEXT[] DEFAULT '{}',
    root_cause TEXT[] DEFAULT '{}',
    resolution TEXT[] DEFAULT '{}',
    prevention TEXT[] DEFAULT '{}',
    commands_used TEXT[] DEFAULT '{}',
    tags TEXT[] DEFAULT '{}',
    severity VARCHAR(10) DEFAULT 'medium',
    environment TEXT,
    services_affected TEXT[] DEFAULT '{}',
    lessons_learned TEXT,
    raw_notes TEXT,
    embedding VECTOR(768),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS incidents;
