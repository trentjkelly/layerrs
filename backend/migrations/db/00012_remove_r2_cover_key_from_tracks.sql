-- +goose Up
-- +goose StatementBegin
ALTER TABLE track DROP COLUMN r2_cover_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE track ADD COLUMN r2_cover_key VARCHAR(64);
-- +goose StatementEnd
