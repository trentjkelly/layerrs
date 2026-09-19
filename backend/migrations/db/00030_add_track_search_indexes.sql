-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_track_search_description
    ON track USING GIN (LOWER(description) gin_trgm_ops)
    WHERE is_valid = true;

CREATE INDEX IF NOT EXISTS idx_artist_search_username
    ON artist USING GIN (LOWER(username) gin_trgm_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_artist_search_username;
DROP INDEX IF EXISTS idx_track_search_description;
-- +goose StatementEnd
