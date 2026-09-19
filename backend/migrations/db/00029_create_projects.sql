-- +goose Up
-- +goose StatementBegin
CREATE TYPE project_track_role AS ENUM ('master', 'stem');

CREATE TABLE project (
    id INT PRIMARY KEY REFERENCES track(id) ON DELETE CASCADE,
    artist_id INT NOT NULL REFERENCES artist(id),
    description VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    plays BIGINT NOT NULL DEFAULT 0,
    likes BIGINT NOT NULL DEFAULT 0,
    layerrs INT NOT NULL DEFAULT 0,
    price_in_cents INT NOT NULL DEFAULT 0,
    is_valid BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX idx_project_artist_id ON project(artist_id);
CREATE INDEX idx_project_created_at ON project(created_at DESC);

CREATE TABLE project_track (
    project_id INT NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    track_id INT NOT NULL REFERENCES track(id) ON DELETE CASCADE,
    role project_track_role NOT NULL,
    position SMALLINT NOT NULL DEFAULT 0,
    PRIMARY KEY (project_id, track_id)
);
CREATE UNIQUE INDEX idx_project_track_one_master ON project_track(project_id) WHERE role = 'master';
CREATE INDEX idx_project_track_track_id ON project_track(track_id);

CREATE TABLE project_tree (
    root_project_id INT NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    child_project_id INT NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    tag derivation_tag NOT NULL DEFAULT 'layerr',
    PRIMARY KEY (root_project_id, child_project_id)
);
CREATE INDEX idx_project_tree_child_id ON project_tree(child_project_id);

INSERT INTO project (id, artist_id, description, created_at, plays, likes, layerrs, price_in_cents, is_valid)
SELECT id, artist_id, description, created_at, plays, likes, layerrs, price_in_cents, is_valid FROM track;
INSERT INTO project_track (project_id, track_id, role)
SELECT id, id, 'master'::project_track_role FROM track;
INSERT INTO project_tree (root_project_id, child_project_id, tag)
SELECT root_id, child_id, tag FROM track_tree;

ALTER TABLE track_stems ADD COLUMN asset_track_id INT REFERENCES track(id) ON DELETE SET NULL;
INSERT INTO track (description, artist_id, wav_r2_track_key, is_valid)
SELECT 'Legacy stem ' || ts.id, p.artist_id, ts.r2_key, TRUE
FROM track_stems ts JOIN project p ON p.id = ts.track_id;
UPDATE track_stems ts
SET asset_track_id = t.id
FROM track t
WHERE t.description = 'Legacy stem ' || ts.id AND t.artist_id = (SELECT artist_id FROM project WHERE id = ts.track_id);
INSERT INTO project_track (project_id, track_id, role, position)
SELECT ts.track_id, ts.asset_track_id, 'stem'::project_track_role, ts.id
FROM track_stems ts WHERE ts.asset_track_id IS NOT NULL;
DROP TABLE track_stems;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS project_tree;
DROP TABLE IF EXISTS project_track;
DROP TABLE IF EXISTS project;
DROP TYPE IF EXISTS project_track_role;
-- +goose StatementEnd
