-- +goose Up
-- +goose StatementBegin
CREATE TYPE artist_stripe_account_status AS ENUM ('pending', 'active', 'inactive', '');

ALTER TABLE artist
ADD COLUMN stripe_account_status artist_stripe_account_status DEFAULT 'inactive';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE artist
DROP COLUMN stripe_account_status;

DROP TYPE artist_stripe_account_status;
-- +goose StatementEnd
