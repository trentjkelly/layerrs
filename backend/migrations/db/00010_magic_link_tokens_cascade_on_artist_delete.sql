-- +goose Up
-- +goose StatementBegin
ALTER TABLE magic_link_tokens
    DROP CONSTRAINT magic_link_tokens_artist_id_fkey,
    ADD CONSTRAINT magic_link_tokens_artist_id_fkey
        FOREIGN KEY (artist_id) REFERENCES artist(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE magic_link_tokens
    DROP CONSTRAINT magic_link_tokens_artist_id_fkey,
    ADD CONSTRAINT magic_link_tokens_artist_id_fkey
        FOREIGN KEY (artist_id) REFERENCES artist(id);
-- +goose StatementEnd
