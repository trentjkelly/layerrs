-- +goose Up
-- +goose StatementBegin
ALTER TABLE track
RENAME COLUMN flac_r2_track_key TO wav_r2_track_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE track
RENAME COLUMN wav_r2_track_key TO flac_r2_track_key;
-- +goose StatementEnd
