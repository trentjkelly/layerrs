-- +goose Up
-- +goose StatementBegin
ALTER TABLE artist ALTER COLUMN username TYPE VARCHAR(30);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE artist ALTER COLUMN username TYPE VARCHAR(64);
-- +goose StatementEnd
