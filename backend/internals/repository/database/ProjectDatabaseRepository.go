package databaseRepository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerrs/internals/entities"
)

type ProjectDatabaseRepository struct{ db *pgxpool.Pool }

func NewProjectDatabaseRepository(db *pgxpool.Pool) *ProjectDatabaseRepository {
	return &ProjectDatabaseRepository{db: db}
}

func (r *ProjectDatabaseRepository) CreateProject(ctx context.Context, project *entities.Project) error {
	query := `INSERT INTO project (id, artist_id, name, description, plays, likes, layerrs, price_in_cents, is_valid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING created_at`
	if err := r.db.QueryRow(ctx, query, project.Id, project.ArtistId, project.Name, project.Description, project.Plays, project.Likes, project.Layerrs, project.PriceInCents, project.IsValid).Scan(&project.CreatedAt); err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	return nil
}

func (r *ProjectDatabaseRepository) UpdateProjectValidity(ctx context.Context, projectId int, isValid bool) error {
	_, err := r.db.Exec(ctx, `UPDATE project SET is_valid = $2 WHERE id = $1`, projectId, isValid)
	if err != nil {
		return fmt.Errorf("failed to update project validity: %w", err)
	}
	return nil
}

func (r *ProjectDatabaseRepository) AddTrack(ctx context.Context, projectId, trackId int, role entities.ProjectTrackRole, position int, name string) error {
	_, err := r.db.Exec(ctx, `INSERT INTO project_track (project_id, track_id, role, position, name) VALUES ($1, $2, $3, $4, $5)`, projectId, trackId, role, position, name)
	if err != nil {
		return fmt.Errorf("failed to attach track to project: %w", err)
	}
	return nil
}

func (r *ProjectDatabaseRepository) CreateDerivation(ctx context.Context, parentProjectId, childProjectId int) error {
	_, err := r.db.Exec(ctx, `INSERT INTO project_tree (root_project_id, child_project_id, tag) VALUES ($1, $2, 'layerr')`, parentProjectId, childProjectId)
	if err != nil {
		return fmt.Errorf("failed to create project derivation: %w", err)
	}
	return nil
}

func (r *ProjectDatabaseRepository) CanUseSourceTracks(ctx context.Context, artistId int, trackIds []int) (bool, error) {
	if len(trackIds) == 0 {
		return true, nil
	}
	var count int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(DISTINCT pt.track_id)
		FROM project_track pt
		JOIN layerrs l ON l.track_id = pt.project_id
		WHERE l.artist_id = $1 AND pt.track_id = ANY($2)
	`, artistId, trackIds).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to authorize source tracks: %w", err)
	}
	return count == len(trackIds), nil
}

func (r *ProjectDatabaseRepository) GetDetail(ctx context.Context, id int) (*entities.ProjectDetail, error) {
	var p entities.ProjectDetail
	err := r.db.QueryRow(ctx, `SELECT p.id, p.name, p.description, a.username, parent.id, parent.name
 FROM project p JOIN artist a ON a.id=p.artist_id
 LEFT JOIN project_tree tree ON tree.child_project_id=p.id
 LEFT JOIN project parent ON parent.id=tree.root_project_id
 WHERE p.id=$1 AND p.is_valid`, id).Scan(&p.Id, &p.Name, &p.Description, &p.ArtistName, &p.ParentId, &p.ParentName)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, `SELECT t.id, COALESCE(NULLIF(pt.name, ''), t.description), pt.role, a.username, COALESCE(NULLIF(t.aac_r2_track_key, ''), t.wav_r2_track_key, ''), COALESCE(t.aac_r2_track_key, '') = ''
 FROM project_track pt JOIN track t ON t.id=pt.track_id JOIN artist a ON a.id=t.artist_id
 WHERE pt.project_id=$1 AND t.is_valid ORDER BY CASE WHEN pt.role='master' THEN 0 ELSE 1 END, pt.position`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	p.Tracks = []entities.ProjectAudio{}
	for rows.Next() {
		var t entities.ProjectAudio
		if err := rows.Scan(&t.Id, &t.Name, &t.Role, &t.ArtistName, &t.Key, &t.Legacy); err != nil {
			return nil, err
		}
		p.Tracks = append(p.Tracks, t)
	}
	return &p, rows.Err()
}
