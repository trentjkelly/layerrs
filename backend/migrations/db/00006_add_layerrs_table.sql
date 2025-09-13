-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS layerrs (
    id SERIAL PRIMARY KEY,
    artist_id INT NOT NULL,
    track_id INT NOT NULL,
    last_layerr_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (artist_id) REFERENCES artist(id),
    FOREIGN KEY (track_id) REFERENCES track(id),
    UNIQUE (artist_id, track_id)
);

CREATE INDEX IF NOT EXISTS idx_layerrs_artist_last_layerr_at ON layerrs(artist_id, last_layerr_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
