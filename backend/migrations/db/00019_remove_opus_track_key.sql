-- +goose Up
-- +goose StatementBegin
ALTER TABLE track
DROP COLUMN opus_r2_track_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE track
ADD COLUMN opus_r2_track_key VARCHAR(255);
-- +goose StatementEnd
