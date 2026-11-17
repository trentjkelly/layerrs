-- +goose Up
ALTER TABLE project ADD COLUMN name VARCHAR(100);
UPDATE project SET name = description;
ALTER TABLE project ALTER COLUMN name SET NOT NULL;
ALTER TABLE project ALTER COLUMN description TYPE VARCHAR(2000);
UPDATE project SET description = '';
ALTER TABLE project_track ADD COLUMN name VARCHAR(255);
UPDATE project_track pt SET name = t.description FROM track t WHERE t.id = pt.track_id;

-- +goose Down
ALTER TABLE project_track DROP COLUMN name;
UPDATE project SET description = name;
ALTER TABLE project ALTER COLUMN description TYPE VARCHAR(100);
ALTER TABLE project DROP COLUMN name;
