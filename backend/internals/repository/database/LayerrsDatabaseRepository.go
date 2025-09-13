package databaseRepository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerrs/internals/entities"
)

type LayerrsDatabaseRepository struct {
	db	*pgxpool.Pool
}

// Constructor for LayerrsDatabaseRepository
func NewLayerrsDatabaseRepository(db *pgxpool.Pool) *LayerrsDatabaseRepository {
	layerrsDatabaseRepository := new(LayerrsDatabaseRepository)
	layerrsDatabaseRepository.db = db
	return layerrsDatabaseRepository
}

// Closes the database pool connection
func (r *LayerrsDatabaseRepository) CloseDB() {
	r.db.Close()
}

// Creates a layerr for a track, or updates the last layerr at date if it already exists
func (r *LayerrsDatabaseRepository) CreateLayerr(ctx context.Context, layerr *entities.Layerr) error {
	query := `INSERT INTO layerrs (artist_id, track_id) VALUES ($1, $2) ON CONFLICT (artist_id, track_id) DO UPDATE SET last_layerr_at = CURRENT_TIMESTAMP RETURNING id;`
	row := r.db.QueryRow(ctx, query, layerr.ArtistId, layerr.TrackId)
	err := row.Scan(&layerr.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *LayerrsDatabaseRepository) ReadLayerr(ctx context.Context, layerr *entities.Layerr) error {
	query := `SELECT id, artist_id, track_id, last_layerr_at FROM layerrs WHERE artist_id=$1 AND track_id=$2 ORDER BY last_layerr_at DESC;`
	row := r.db.QueryRow(ctx, query, layerr.ArtistId, layerr.TrackId)
	err := row.Scan(&layerr.Id, &layerr.ArtistId, &layerr.TrackId, &layerr.LastLayerrAt)
	if err != nil {
		return err
	}
	return nil
}