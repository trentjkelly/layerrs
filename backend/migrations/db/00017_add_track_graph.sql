-- +goose Up
-- +goose StatementBegin
CREATE TABLE graphs (
    id           SERIAL PRIMARY KEY,
    total_tracks INT NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE track_graphs (
    track_id   INT PRIMARY KEY REFERENCES track(id),
    graph_id   INT NOT NULL REFERENCES graphs(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE track_graphs;
DROP TABLE graphs;
-- +goose StatementEnd
