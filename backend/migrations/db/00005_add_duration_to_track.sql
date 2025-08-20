-- +goose Up
-- +goose StatementBegin
ALTER TABLE track 
ADD COLUMN duration NUMERIC(8,3)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
