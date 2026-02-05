-- +goose Up
-- +goose StatementBegin
ALTER TABLE artist DROP COLUMN password;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
