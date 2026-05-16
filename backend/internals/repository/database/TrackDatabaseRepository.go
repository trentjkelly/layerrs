package databaseRepository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerrs/internals/entities"
)

type TrackDatabaseRepository struct {
	db *pgxpool.Pool
}

// Constructor for TrackDatabaseRepository
func NewTrackDatabaseRepository(db *pgxpool.Pool) *TrackDatabaseRepository {
	trackDatabaseRepository := new(TrackDatabaseRepository)
	trackDatabaseRepository.db = db
	return trackDatabaseRepository
}

// Closes the database pool connection
func (r *TrackDatabaseRepository) CloseDB() {
	r.db.Close()
}

// Adds a Track to the database, but only the non-optional fields
func (r *TrackDatabaseRepository) CreateTrack(ctx context.Context, track *entities.Track) error {
	query := `INSERT INTO track (description, artist_id) VALUES ($1, $2) RETURNING id`
	row := r.db.QueryRow(ctx, query, track.Description, track.ArtistId)

	err := row.Scan(&track.Id)
	if err != nil {
		return fmt.Errorf("failed to scan rows in CreateTrack: %w", err)
	}

	return nil
}

// Gets a Track from the database based on their id
func (r *TrackDatabaseRepository) ReadTrackById(ctx context.Context, track *entities.Track) error {
	query := `SELECT id, description, artist_id, wav_r2_track_key, aac_r2_track_key, created_at, plays, likes, layerrs, is_valid, duration, price_in_cents FROM track WHERE id=$1;`
	row := r.db.QueryRow(ctx, query, track.Id)

	// Potential NULL Values
	var wavR2TrackKey sql.NullString
	var aacR2TrackKey sql.NullString

	err := row.Scan(&track.Id, &track.Description, &track.ArtistId, &wavR2TrackKey, &aacR2TrackKey, &track.CreatedAt, &track.Plays, &track.Likes, &track.Layerrs, &track.IsValid, &track.TrackDuration, &track.PriceInCents)
	if err != nil {
		return fmt.Errorf("failed to scan rows in ReadTrackByID: %w", err)
	}

	// Potential NULL values converted to empty strings
	if wavR2TrackKey.Valid {
		track.WavR2TrackKey = wavR2TrackKey.String
	} else {
		track.WavR2TrackKey = ""
	}

	if aacR2TrackKey.Valid {
		track.AacR2TrackKey = aacR2TrackKey.String
	} else {
		track.AacR2TrackKey = ""
	}

	return nil
}

// Updates the information for a Track in the database
func (r *TrackDatabaseRepository) UpdateTrack(ctx context.Context, track *entities.Track) error {
	query := `UPDATE track SET description=$2, wav_r2_track_key=$3, aac_r2_track_key=$4, is_valid=$5, duration=$6 WHERE id=$1 RETURNING description;`
	row := r.db.QueryRow(ctx, query, track.Id, track.Description, track.WavR2TrackKey, track.AacR2TrackKey, track.IsValid, track.TrackDuration)

	err := row.Scan(&track.Description)
	if err != nil {
		return fmt.Errorf("failed to scan rows in UpdateTrack: %w", err)
	}

	return nil
}

// Deletes a Track from the database given the trackId
func (r *TrackDatabaseRepository) DeleteTrack(ctx context.Context, track *entities.Track) error {
	query := `DELETE FROM track WHERE id=$1 RETURNING id;`
	row := r.db.QueryRow(ctx, query, track.Id)

	err := row.Scan(&track.Id)
	if err != nil {
		return fmt.Errorf("failed to scan rows in DeleteTrack: %w", err)
	}

	return nil
}

// Updates the price of a track
func (r *TrackDatabaseRepository) UpdateTrackPrice(ctx context.Context, trackId int, priceInCents int) error {
	query := `UPDATE track SET updated_at=$2, price_in_cents=$3 WHERE id=$1 RETURNING updated_at;`
	row := r.db.QueryRow(ctx, query, trackId, time.Now(), priceInCents)

	var updatedAt time.Time
	err := row.Scan(&updatedAt)
	if err != nil {
		return fmt.Errorf("failed to update track price: %w", err)
	}

	return nil
}

// Increases the number of plays on a Track when a user plays a song
func (r *TrackDatabaseRepository) IncrementPlays(ctx context.Context, track *entities.Track) error {
	query := `UPDATE track SET plays = plays + 1 WHERE id=$1 RETURNING plays;`
	row := r.db.QueryRow(ctx, query, track.Id)

	err := row.Scan(&track.Plays)
	if err != nil {
		return fmt.Errorf("failed to scan rows in IncrementPlays: %w", err)
	}

	return nil
}

// Increases the number of likes on a Track when a user likes it
func (r *TrackDatabaseRepository) IncrementLikes(ctx context.Context, track *entities.Track) error {
	query := `UPDATE track SET likes = likes + 1 WHERE id=$1 RETURNING likes;`
	row := r.db.QueryRow(ctx, query, track.Id)

	err := row.Scan(&track.Likes)
	if err != nil {
		return fmt.Errorf("failed to scan rows in IncrementLikes: %w", err)
	}

	return nil
}

// Decreases the number of likes on a Track when a user likes it
func (r *TrackDatabaseRepository) DecrementLikes(ctx context.Context, track *entities.Track) error {
	query := `UPDATE track SET likes = likes - 1 WHERE id=$1 RETURNING likes;`
	row := r.db.QueryRow(ctx, query, track.Id)

	err := row.Scan(&track.Likes)
	if err != nil {
		return fmt.Errorf("failed to scan rows in DecrementLikes: %w", err)
	}

	return nil
}

// Gets a single track's recommendation data by its ID
func (r *TrackDatabaseRepository) ReadOneTrackById(ctx context.Context, trackId int, artistId int) (entities.TrackInfo, error) {
	query := `
		SELECT t.id, t.description, t.artist_id, a.username, a.r2_image_key, t.likes, t.layerrs, t.duration, w.waveform_data, t.price_in_cents,
		CASE WHEN alt.artist_id IS NOT NULL THEN true ELSE false END as is_liked
		FROM track t
		JOIN artist a ON t.artist_id = a.id
		LEFT JOIN waveform w ON w.track_id = t.id
		LEFT JOIN artist_likes_track alt ON alt.track_id = t.id AND alt.artist_id = $2
		WHERE t.id = $1 AND t.is_valid = true;
	`

	var rec entities.TrackInfo
	var waveformData []int
	var r2ImageKey sql.NullString
	err := r.db.QueryRow(ctx, query, trackId, artistId).Scan(
		&rec.Id, &rec.Description, &rec.ArtistId, &rec.ArtistName, &r2ImageKey,
		&rec.Likes, &rec.Layerrs, &rec.Duration, &waveformData, &rec.PriceInCents, &rec.IsLiked,
	)
	if err != nil {
		return rec, fmt.Errorf("failed to scan row in ReadOneTrackById: %w", err)
	}
	if r2ImageKey.Valid {
		rec.R2ImageKey = r2ImageKey.String
	}
	if waveformData != nil {
		rec.WaveformData = waveformData
	} else {
		rec.WaveformData = []int{}
	}
	return rec, nil
}

// Gets the top N most recent tracks with full track data -- used for recommendations algorithm
func (r *TrackDatabaseRepository) ReadNTracksByDate(ctx context.Context, offset int, artistId int) ([]entities.TrackInfo, error) {
	query := `
		SELECT t.id, t.description, t.artist_id, a.username, a.r2_image_key, t.likes, t.layerrs, t.duration, w.waveform_data, t.price_in_cents,
		CASE WHEN alt.artist_id IS NOT NULL THEN true ELSE false END as is_liked
		FROM track t
		JOIN artist a ON t.artist_id = a.id
		LEFT JOIN waveform w ON w.track_id = t.id
		LEFT JOIN artist_likes_track alt ON alt.track_id = t.id AND alt.artist_id = $2
		WHERE t.is_valid = true
		ORDER BY t.created_at DESC
		LIMIT 8 OFFSET $1;
	`

	rows, err := r.db.Query(ctx, query, offset, artistId)
	if err != nil {
		return nil, fmt.Errorf("failed to query rows in ReadNTracksByDate: %w", err)
	}
	defer rows.Close()

	var recs []entities.TrackInfo

	for rows.Next() {
		var rec entities.TrackInfo
		var waveformData []int
		var r2ImageKey sql.NullString
		err = rows.Scan(&rec.Id, &rec.Description, &rec.ArtistId, &rec.ArtistName, &r2ImageKey, &rec.Likes, &rec.Layerrs, &rec.Duration, &waveformData, &rec.PriceInCents, &rec.IsLiked)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows in ReadNTracksByDate: %w", err)
		}
		if r2ImageKey.Valid {
			rec.R2ImageKey = r2ImageKey.String
		}
		if waveformData != nil {
			rec.WaveformData = waveformData
		} else {
			rec.WaveformData = []int{}
		}
		recs = append(recs, rec)
	}

	return recs, nil
}

// Gets full track info for a list of track IDs
func (r *TrackDatabaseRepository) ReadTracksByIds(ctx context.Context, trackIds []int, artistId int) ([]entities.TrackInfo, error) {
	query := `
		SELECT t.id, t.description, t.artist_id, a.username, a.r2_image_key, t.likes, t.layerrs, t.duration, w.waveform_data, t.price_in_cents,
		CASE WHEN alt.artist_id IS NOT NULL THEN true ELSE false END as is_liked
		FROM track t
		JOIN artist a ON t.artist_id = a.id
		LEFT JOIN waveform w ON w.track_id = t.id
		LEFT JOIN artist_likes_track alt ON alt.track_id = t.id AND alt.artist_id = $2
		WHERE t.is_valid = true AND t.id = ANY($1);
	`

	rows, err := r.db.Query(ctx, query, trackIds, artistId)
	if err != nil {
		return nil, fmt.Errorf("failed to query rows in ReadTracksByIds: %w", err)
	}
	defer rows.Close()

	var recs []entities.TrackInfo

	for rows.Next() {
		var rec entities.TrackInfo
		var waveformData []int
		var r2ImageKey sql.NullString
		err = rows.Scan(&rec.Id, &rec.Description, &rec.ArtistId, &rec.ArtistName, &r2ImageKey, &rec.Likes, &rec.Layerrs, &rec.Duration, &waveformData, &rec.PriceInCents, &rec.IsLiked)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows in ReadTracksByIds: %w", err)
		}
		if r2ImageKey.Valid {
			rec.R2ImageKey = r2ImageKey.String
		}
		if waveformData != nil {
			rec.WaveformData = waveformData
		} else {
			rec.WaveformData = []int{}
		}
		recs = append(recs, rec)
	}

	return recs, nil
}

// Gets the top N tracks by likes with full track data -- used for recommendations algorithm
func (r *TrackDatabaseRepository) ReadNTracksByLikes(ctx context.Context, offset int) ([]entities.TrackInfo, error) {
	query := `
		SELECT t.id, t.description, t.artist_id, a.username, a.r2_image_key, t.likes, t.layerrs, t.duration, w.waveform_data, t.price_in_cents
		FROM track t
		JOIN artist a ON t.artist_id = a.id
		LEFT JOIN waveform w ON w.track_id = t.id
		WHERE t.is_valid = true
		ORDER BY t.likes DESC
		LIMIT 8 OFFSET $1;
	`

	rows, err := r.db.Query(ctx, query, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query rows in ReadNTracksByLikes: %w", err)
	}
	defer rows.Close()

	var recs []entities.TrackInfo

	for rows.Next() {
		var rec entities.TrackInfo
		var waveformData []int
		var r2ImageKey sql.NullString
		err = rows.Scan(&rec.Id, &rec.Description, &rec.ArtistId, &rec.ArtistName, &r2ImageKey, &rec.Likes, &rec.Layerrs, &rec.Duration, &waveformData, &rec.PriceInCents)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows in ReadNTracksByLikes: %w", err)
		}
		if r2ImageKey.Valid {
			rec.R2ImageKey = r2ImageKey.String
		}
		if waveformData != nil {
			rec.WaveformData = waveformData
		} else {
			rec.WaveformData = []int{}
		}
		recs = append(recs, rec)
	}

	return recs, nil
}
