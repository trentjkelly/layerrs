-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS magic_link_tokens (
    id SERIAL PRIMARY KEY,
    hashed_token VARCHAR(255) NOT NULL,
    artist_id INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (artist_id) REFERENCES artist(id)
);

CREATE INDEX IF NOT EXISTS idx_magic_link_tokens_token ON magic_link_tokens(hashed_token);
CREATE INDEX IF NOT EXISTS idx_magic_link_tokens_artist_id ON magic_link_tokens(artist_id);
CREATE INDEX IF NOT EXISTS idx_magic_link_tokens_created_at ON magic_link_tokens(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
