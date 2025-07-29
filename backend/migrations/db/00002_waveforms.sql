-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS waveform (
    id SERIAL PRIMARY KEY,
    track_id INT NOT NULL,
    waveform_data INT[] NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (track_id) REFERENCES track(id)
);

CREATE INDEX IF NOT EXISTS idx_waveform_track_id ON waveform(track_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
