-- +goose Up
ALTER TABLE issues 
ADD COLUMN steps TEXT,
ADD COLUMN tags TEXT[];

-- +goose Down
ALTER TABLE issues 
DROP COLUMN steps,
DROP COLUMN tags;