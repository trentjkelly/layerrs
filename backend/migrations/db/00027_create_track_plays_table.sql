-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS track_plays (
    id BIGSERIAL PRIMARY KEY,
    track_id INT NOT NULL REFERENCES track(id) ON DELETE CASCADE,
    artist_id INT NOT NULL REFERENCES artist(id) ON DELETE CASCADE,
    played_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_track_plays_track_artist ON track_plays(track_id, artist_id);
CREATE INDEX IF NOT EXISTS idx_track_plays_played_at ON track_plays(played_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS track_plays;
-- +goose StatementEnd
