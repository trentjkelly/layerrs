package databaseRepository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerrs/internals/entities"
)

type WaveformDatabaseRepository struct {
	db *pgxpool.Pool
}

// Constructor for a new waveform database repository
func NewWaveformDatabaseRepository(db *pgxpool.Pool) *WaveformDatabaseRepository {
	waveformDbRepo := new(WaveformDatabaseRepository)
	waveformDbRepo.db = db

	return waveformDbRepo
}

// Creates a waveform database entry
func (r *WaveformDatabaseRepository) CreateWaveform(ctx context.Context, waveform *entities.Waveform) error {
	query := `INSERT INTO waveform (track_id, waveform_data) VALUES ($1, $2) RETURNING id;`
	row := r.db.QueryRow(ctx, query, waveform.TrackId, waveform.WaveformData)

	err := row.Scan(&waveform.Id)
	if err != nil {
		return fmt.Errorf("failed to create waveform in the db: %w", err)
	}

	return nil
}

// Gets a waveform from the database based on the track id
func (r *WaveformDatabaseRepository) GetWaveform(ctx context.Context, waveform *entities.Waveform) (error) {
	query := `SELECT id, waveform_data FROM waveform WHERE track_id = $1;`
	row := r.db.QueryRow(ctx, query, waveform.TrackId)

	err := row.Scan(&waveform.Id, &waveform.WaveformData)
	if err != nil {
		return fmt.Errorf("failed to get waveform from the db: %w", err)
	}

	return nil
}