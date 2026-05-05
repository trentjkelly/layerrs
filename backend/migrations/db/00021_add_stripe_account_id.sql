-- +goose Up
-- +goose StatementBegin
ALTER TABLE artist
ADD COLUMN stripe_account_id VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE artist
DROP COLUMN stripe_account_id;
-- +goose StatementEnd
