-- +goose Up
-- +goose StatementBegin
ALTER TABLE track 
ADD COLUMN flac_r2_track_key VARCHAR(255);

ALTER TABLE track
ADD COLUMN opus_r2_track_key VARCHAR(255);

ALTER TABLE track
ADD COLUMN aac_r2_track_key VARCHAR(255);

ALTER TABLE track
DROP COLUMN r2_track_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
