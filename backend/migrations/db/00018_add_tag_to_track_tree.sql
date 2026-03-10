-- +goose Up
-- +goose StatementBegin
CREATE TYPE derivation_tag AS ENUM ('layerr', 'stem');

ALTER TABLE track_tree
    ADD COLUMN tag derivation_tag NOT NULL DEFAULT 'layerr';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE track_tree DROP COLUMN tag;

DROP TYPE derivation_tag;
-- +goose StatementEnd
