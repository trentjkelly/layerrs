-- +goose Up
-- +goose StatementBegin
ALTER TABLE track
ADD COLUMN is_valid BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
