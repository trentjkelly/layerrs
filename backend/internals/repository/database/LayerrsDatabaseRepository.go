package databaseRepository

import (
	"context"
	"fmt"
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

// Reads all layerrs for an artist with full track and artist info, sorted by last layerr at date
func (r *LayerrsDatabaseRepository) ReadLayerrsWithTracks(ctx context.Context, artistId int) ([]entities.LayerrTrack, error) {
	query := `
		SELECT t.id, t.description, t.artist_id, a.username, a.r2_image_key,
		       t.likes, t.layerrs, t.duration, w.waveform_data,
		       CASE WHEN alt.artist_id IS NOT NULL THEN true ELSE false END as is_liked,
		       t.color, l.last_layerr_at
		FROM layerrs l
		JOIN track t ON l.track_id = t.id
		JOIN artist a ON t.artist_id = a.id
		LEFT JOIN waveform w ON w.track_id = t.id
		LEFT JOIN artist_likes_track alt ON alt.track_id = t.id AND alt.artist_id = $1
		WHERE l.artist_id = $1
		ORDER BY l.last_layerr_at DESC;
	`

	rows, err := r.db.Query(ctx, query, artistId)
	if err != nil {
		return nil, fmt.Errorf("failed to query layerrs with tracks from database: %w", err)
	}
	defer rows.Close()

	var layerrTracks []entities.LayerrTrack
	for rows.Next() {
		var lt entities.LayerrTrack
		var waveformData []int
		err = rows.Scan(
			&lt.Id, &lt.Description, &lt.ArtistId, &lt.ArtistName, &lt.R2ImageKey,
			&lt.Likes, &lt.Layerrs, &lt.Duration, &waveformData, &lt.IsLiked,
			&lt.Color, &lt.LastLayerrAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan layerrs with tracks from database: %w", err)
		}
		if waveformData != nil {
			lt.WaveformData = waveformData
		} else {
			lt.WaveformData = []int{}
		}
		layerrTracks = append(layerrTracks, lt)
	}
	return layerrTracks, nil
}

// Reads all layerrs for an artist sorted by last layerr at date
func (r *LayerrsDatabaseRepository) ReadLayerrs(ctx context.Context, artistId int) ([]*entities.Layerr, error) {
	query := `SELECT id, artist_id, track_id, last_layerr_at FROM layerrs WHERE artist_id=$1 ORDER BY last_layerr_at DESC;`
	rows, err := r.db.Query(ctx, query, artistId)
	if err != nil {
		return nil, fmt.Errorf("failed to query layerrs from database: %w", err)
	}
	
	var layerrs []*entities.Layerr
	for rows.Next() {
		layerr := new(entities.Layerr)

		err = rows.Scan(&layerr.Id, &layerr.ArtistId, &layerr.TrackId, &layerr.LastLayerrAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan layerrs from database: %w", err)
		}

		layerrs = append(layerrs, layerr)
	}
	return layerrs, nil
}