-- +goose Up
-- +goose StatementBegin
ALTER TABLE track RENAME COLUMN name TO description;
ALTER TABLE track ALTER COLUMN description TYPE VARCHAR(100);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE track ALTER COLUMN description TYPE VARCHAR(64);
ALTER TABLE track RENAME COLUMN description TO name;
-- +goose StatementEnd
