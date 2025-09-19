-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS genre (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
)

CREATE TABLE IF NOT EXISTS genre_mods (
    id SERIAL PRIMARY KEY,
    genre_id INT,
    artist_id INT,
    is_founder BOOLEAN,
    is_active BOOLEAN,
    UNIQUE (genre_id, artist_id),
    FOREIGN KEY genre_id REFERENCES genre(id),
    FOREIGN KEY artist_id REFERENCES artist(id)
)

CREATE TABLE IF NOT EXISTS genre_tracks (
    id SERIAL PRIMARY KEY,
    genre_id INT,
    track_id INT,
    added_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(genre_id, artist_id),
    FOREIGN KEY genre_id REFERENCES genre(id),
    FOREIGN KEY track_id REFERENCES track(id)
)

-- TODO: Add indexes for frequent queries after determining database reads/writes

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
