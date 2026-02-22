-- +goose Up
-- +goose StatementBegin
ALTER TABLE artist ALTER COLUMN username DROP NOT NULL;
ALTER TABLE artist DROP COLUMN name;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE artist ADD COLUMN name VARCHAR(64) NOT NULL DEFAULT '';
ALTER TABLE artist ALTER COLUMN username SET NOT NULL;
-- +goose StatementEnd
