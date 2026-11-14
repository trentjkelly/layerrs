-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS track_stems (
    id BIGSERIAL PRIMARY KEY,
    track_id INT NOT NULL REFERENCES track(id) ON DELETE CASCADE,
    original_filename VARCHAR(255) NOT NULL,
    r2_key VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_track_stems_track_id ON track_stems(track_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS track_stems;
-- +goose StatementEnd
