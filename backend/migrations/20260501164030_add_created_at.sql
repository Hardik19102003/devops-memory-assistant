-- +goose Up
ALTER TABLE issues 
ADD COLUMN created_at TIMESTAMP DEFAULT NOW();

UPDATE issues 
SET created_at = NOW() 
WHERE created_at IS NULL;

ALTER TABLE issues 
ALTER COLUMN created_at SET NOT NULL;

-- +goose Down
ALTER TABLE issues 
DROP COLUMN created_at;