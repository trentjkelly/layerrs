-- +goose Up
-- +goose StatementBegin
ALTER TABLE artist ADD COLUMN can_post BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE artist DROP COLUMN can_post;
-- +goose StatementEnd
