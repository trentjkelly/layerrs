-- +goose Up
-- +goose StatementBegin

CREATE TABLE page (
    id SERIAL PRIMARY KEY,
    editor_id INTEGER REFERENCES artist(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (editor_id, name)
);

CREATE INDEX idx_page_editor_id ON page(editor_id);

CREATE TABLE page_track (
    id SERIAL PRIMARY KEY,
    page_id INTEGER NOT NULL REFERENCES page(id) ON DELETE CASCADE,
    track_id INTEGER NOT NULL REFERENCES track(id) ON DELETE CASCADE,
    added_at TIMESTAMPTZ DEFAULT NOW(),
    recommender_id INTEGER REFERENCES artist(id) ON DELETE SET NULL,
    UNIQUE (page_id, track_id)
);

CREATE INDEX idx_page_track_page_id_added_at ON page_track(page_id, added_at DESC);
CREATE INDEX idx_page_track_track_id ON page_track(track_id);

CREATE TABLE page_follower (
    id SERIAL PRIMARY KEY,
    page_id INTEGER NOT NULL REFERENCES page(id) ON DELETE CASCADE,
    artist_id INTEGER NOT NULL REFERENCES artist(id) ON DELETE CASCADE,
    followed_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (page_id, artist_id)
);

CREATE INDEX idx_page_follower_page_id ON page_follower(page_id);
CREATE INDEX idx_page_follower_artist_id ON page_follower(artist_id);

CREATE TABLE page_submission (
    id SERIAL PRIMARY KEY,
    page_id INTEGER NOT NULL REFERENCES page(id) ON DELETE CASCADE,
    track_id INTEGER NOT NULL REFERENCES track(id) ON DELETE CASCADE,
    submitter_id INTEGER NOT NULL REFERENCES artist(id) ON DELETE CASCADE,
    note TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (page_id, track_id)
);

CREATE INDEX idx_page_submission_page_id ON page_submission(page_id);
CREATE INDEX idx_page_submission_track_id ON page_submission(track_id);
CREATE INDEX idx_page_submission_submitter_id ON page_submission(submitter_id);

CREATE TABLE page_track_note (
    id SERIAL PRIMARY KEY,
    page_track_id INTEGER NOT NULL REFERENCES page_track(id) ON DELETE CASCADE,
    note TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_page_track_note_page_track_id ON page_track_note(page_track_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS page_track_note;
DROP TABLE IF EXISTS page_submission;
DROP TABLE IF EXISTS page_follower;
DROP TABLE IF EXISTS page_track;
DROP TABLE IF EXISTS page;

-- +goose StatementEnd
