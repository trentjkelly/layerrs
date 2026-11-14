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
	query := `INSERT INTO project (id, artist_id, description, plays, likes, layerrs, price_in_cents, is_valid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING created_at`
	if err := r.db.QueryRow(ctx, query, project.Id, project.ArtistId, project.Description, project.Plays, project.Likes, project.Layerrs, project.PriceInCents, project.IsValid).Scan(&project.CreatedAt); err != nil {
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

func (r *ProjectDatabaseRepository) AddTrack(ctx context.Context, projectId, trackId int, role entities.ProjectTrackRole, position int) error {
	_, err := r.db.Exec(ctx, `INSERT INTO project_track (project_id, track_id, role, position) VALUES ($1, $2, $3, $4)`, projectId, trackId, role, position)
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
