package databaseRepository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerrs/internals/entities"
)

type TrackPlayDatabaseRepository struct {
	db *pgxpool.Pool
}

// Constructor for TrackPlayDatabaseRepository
func NewTrackPlayDatabaseRepository(db *pgxpool.Pool) *TrackPlayDatabaseRepository {
	trackPlayDatabaseRepository := new(TrackPlayDatabaseRepository)
	trackPlayDatabaseRepository.db = db
	return trackPlayDatabaseRepository
}

// Closes the database pool connection
func (r *TrackPlayDatabaseRepository) CloseDB() {
	r.db.Close()
}

// Records that an artist played a track
func (r *TrackPlayDatabaseRepository) CreateTrackPlay(ctx context.Context, play *entities.TrackPlay) error {
	query := `INSERT INTO track_plays (track_id, artist_id) VALUES ($1, $2) RETURNING id`
	row := r.db.QueryRow(ctx, query, play.TrackId, play.ArtistId)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to create track play: %w", err)
	}

	return nil
}
