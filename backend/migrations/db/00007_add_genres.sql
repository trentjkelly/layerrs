-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS genre (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS genre_mods (
    id SERIAL PRIMARY KEY,
    genre_id INT,
    artist_id INT,
    is_founder BOOLEAN,
    added_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (genre_id, artist_id),
    FOREIGN KEY (genre_id) REFERENCES genre(id),
    FOREIGN KEY (artist_id) REFERENCES artist(id)
);

CREATE TABLE IF NOT EXISTS genre_tracks (
    id SERIAL PRIMARY KEY,
    genre_id INT,
    track_id INT,
    added_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(genre_id, track_id),
    FOREIGN KEY (genre_id) REFERENCES genre(id),
    FOREIGN KEY (track_id) REFERENCES track(id)
);

CREATE INDEX idx_genre_name ON genre(name);
CREATE INDEX idx_genre_mods_genre_id ON genre_mods(genre_id);
CREATE INDEX idx_genre_mods_artist_id ON genre_mods(artist_id);
CREATE INDEX idx_genre_mods_is_founder ON genre_mods(genre_id, is_founder);
CREATE INDEX idx_genre_tracks_genre_id ON genre_tracks(genre_id);
CREATE INDEX idx_genre_tracks_track_id ON genre_tracks(track_id);
CREATE INDEX idx_genre_mods_added_at ON genre_mods(added_at);
CREATE INDEX idx_genre_tracks_added_at ON genre_tracks(added_at);
CREATE INDEX idx_genre_tracks_genre_added ON genre_tracks(genre_id, added_at DESC);
-- TODO: Add indexes for frequent queries after determining database reads/writes

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
