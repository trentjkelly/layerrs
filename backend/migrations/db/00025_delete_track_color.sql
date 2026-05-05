-- +goose Up
-- +goose StatementBegin
ALTER TABLE track DROP COLUMN color;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE track ADD COLUMN color VARCHAR(255);
-- +goose StatementEnd
