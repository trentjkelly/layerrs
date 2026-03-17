-- +goose Up
-- +goose StatementBegin
ALTER TABLE track ADD COLUMN color VARCHAR(10) NOT NULL DEFAULT 'violet';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE track DROP COLUMN color;
-- +goose StatementEnd
