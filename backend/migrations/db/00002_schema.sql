CREATE TABLE artist (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    username VARCHAR(64) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(72) NOT NULL DEFAULT 0,
    bio TEXT,
    r2_image_key VARCHAR(64),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE track (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    artist_id INT,
    r2_track_key VARCHAR(64),
    r2_cover_key VARCHAR(64),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    plays BIGINT DEFAULT 0,
    likes BIGINT DEFAULT 0,
    FOREIGN KEY (artist_id) REFERENCES artist(id)
);

CREATE TABLE artist_likes_track (
    id BIGSERIAL PRIMARY KEY,
    artist_id INT NOT NULL,
    track_id INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (artist_id) REFERENCES artist(id),
    FOREIGN KEY (track_id) REFERENCES track(id),
    UNIQUE (artist_id, track_id)
);

CREATE TABLE track_tree (
    root_id INT NOT NULL,
    child_id INT NOT NULL,
    PRIMARY KEY (root_id, child_id),
    FOREIGN KEY (root_id) REFERENCES track(id),
    FOREIGN KEY (child_id) REFERENCES track(id)
);

-- Indexes
CREATE INDEX idx_artist_id ON track(artist_id);
CREATE INDEX idx_artist_likes_artist_id ON artist_likes_track(artist_id);
CREATE INDEX idx_artist_likes_track_id ON artist_likes_track(track_id);
CREATE INDEX idx_track_tree_root_id ON track_tree(root_id);
CREATE INDEX idx_track_tree_child_id ON track_tree(child_id);