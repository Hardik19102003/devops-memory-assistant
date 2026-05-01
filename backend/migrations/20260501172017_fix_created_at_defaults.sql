-- +goose Up

-- 1. Fill old rows
UPDATE issues 
SET created_at = NOW() 
WHERE created_at IS NULL;

-- 2. Make it NOT NULL
ALTER TABLE issues 
ALTER COLUMN created_at SET NOT NULL;

-- +goose Down

ALTER TABLE issues 
ALTER COLUMN created_at DROP NOT NULL;