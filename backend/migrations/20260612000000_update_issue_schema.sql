-- +goose Up
ALTER TABLE issues 
ADD COLUMN IF NOT EXISTS title TEXT,
ADD COLUMN IF NOT EXISTS summary TEXT,
ADD COLUMN IF NOT EXISTS symptoms TEXT[] DEFAULT '{}',
ADD COLUMN IF NOT EXISTS evidence TEXT[] DEFAULT '{}',
ADD COLUMN IF NOT EXISTS root_cause TEXT[] DEFAULT '{}',
ADD COLUMN IF NOT EXISTS resolution TEXT[] DEFAULT '{}',
ADD COLUMN IF NOT EXISTS prevention TEXT[] DEFAULT '{}',
ADD COLUMN IF NOT EXISTS commands_used TEXT[] DEFAULT '{}',
ADD COLUMN IF NOT EXISTS tags TEXT[] DEFAULT '{}',
ADD COLUMN IF NOT EXISTS severity VARCHAR(10) DEFAULT 'medium',
ADD COLUMN IF NOT EXISTS environment TEXT,
ADD COLUMN IF NOT EXISTS services_affected TEXT[] DEFAULT '{}',
ADD COLUMN IF NOT EXISTS lessons_learned TEXT,
ADD COLUMN IF NOT EXISTS raw_notes TEXT,
ADD COLUMN IF NOT EXISTS embedding VECTOR(768);

-- Rename existing columns to match new schema (if needed, but we'll keep old for migration safety)
-- We'll keep the old columns and let the service layer handle mapping, or we can migrate data.
-- For simplicity, we'll assume we are starting fresh or can migrate data separately.
-- If we want to preserve existing data, we would need to map:
-- error -> title? or summary? We'll leave that to the application.
-- For now, we just add the new columns.

-- +goose Down
ALTER TABLE issues 
DROP COLUMN IF EXISTS title,
DROP COLUMN IF EXISTS summary,
DROP COLUMN IF EXISTS symptoms,
DROP COLUMN IF EXISTS evidence,
DROP COLUMN IF EXISTS root_cause,
DROP COLUMN IF EXISTS resolution,
DROP COLUMN IF EXISTS prevention,
DROP COLUMN IF EXISTS commands_used,
DROP COLUMN IF EXISTS tags,
DROP COLUMN IF EXISTS severity,
DROP COLUMN IF EXISTS environment,
DROP COLUMN IF EXISTS services_affected,
DROP COLUMN IF EXISTS lessons_learned,
DROP COLUMN IF EXISTS raw_notes,
DROP COLUMN IF EXISTS embedding;