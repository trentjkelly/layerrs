-- +goose Up
-- +goose StatementBegin
ALTER TABLE track ADD COLUMN price_in_cents INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE track DROP COLUMN price_in_cents;
-- +goose StatementEnd
