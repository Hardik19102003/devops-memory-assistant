-- +goose Up
ALTER TABLE issues 
ADD COLUMN created_at TIMESTAMP DEFAULT NOW();

-- +goose Down
ALTER TABLE issues 
DROP COLUMN created_at;
