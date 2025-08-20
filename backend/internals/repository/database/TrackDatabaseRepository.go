package databaseRepository

import (
	"database/sql"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerrs/internals/entities"
	"fmt"
)

type TrackDatabaseRepository struct {
	db	*pgxpool.Pool
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
	query := `INSERT INTO track (name, artist_id) VALUES ($1, $2) RETURNING id`
	row := r.db.QueryRow(ctx, query, track.Name, track.ArtistId)
	
	err := row.Scan(&track.Id)
	if err != nil {
		return fmt.Errorf("failed to scan rows in CreateTrack: %w", err)
	}

	return nil
}

// Gets a Track from the database based on their id
func (r *TrackDatabaseRepository) ReadTrackById(ctx context.Context, track *entities.Track) error {
	query := `SELECT id, name, artist_id, flac_r2_track_key, opus_r2_track_key, aac_r2_track_key, r2_cover_key, created_at, plays, likes, layerrs, is_valid, duration FROM track WHERE id=$1;`
	row := r.db.QueryRow(ctx, query, track.Id)
	
	// Potential NULL Values
	var flacR2TrackKey sql.NullString
	var opusR2TrackKey sql.NullString
	var aacR2TrackKey sql.NullString
	var r2CoverKey sql.NullString

	err := row.Scan(&track.Id, &track.Name, &track.ArtistId, &flacR2TrackKey, &opusR2TrackKey, &aacR2TrackKey, &r2CoverKey, &track.CreatedAt, &track.Plays, &track.Likes, &track.Layerrs, &track.IsValid, &track.TrackDuration)
	if err != nil {
		return fmt.Errorf("failed to scan rows in ReadTrackByID: %w", err)
	}

	// Potential NULL values converted to empty strings
	if flacR2TrackKey.Valid {
		track.FlacR2TrackKey = flacR2TrackKey.String
	} else {
		track.FlacR2TrackKey = ""
	}

	if opusR2TrackKey.Valid {
		track.OpusR2TrackKey = opusR2TrackKey.String
	} else {
		track.OpusR2TrackKey = ""
	}

	if aacR2TrackKey.Valid {
		track.AacR2TrackKey = aacR2TrackKey.String
	} else {
		track.AacR2TrackKey = ""
	}

	if r2CoverKey.Valid {
		track.R2CoverKey = r2CoverKey.String
	} else {
		track.R2CoverKey = ""
	}

	return nil
}

// Updates the information for a Track in the database
func (r *TrackDatabaseRepository) UpdateTrack(ctx context.Context, track *entities.Track) error {
	query := `UPDATE track SET name=$2, flac_r2_track_key=$3, opus_r2_track_key=$4, aac_r2_track_key=$5, r2_cover_key=$6, is_valid=$7, duration=$8 WHERE id=$1 RETURNING name;`
	row := r.db.QueryRow(ctx, query, track.Id, track.Name, track.FlacR2TrackKey, track.OpusR2TrackKey, track.AacR2TrackKey, track.R2CoverKey, track.IsValid, track.TrackDuration)
	
	err := row.Scan(&track.Name)
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

// Gets the top N tracks by likes -- used for recommendations algorithm
func (r *TrackDatabaseRepository) ReadNTracksByLikes(ctx context.Context, offset int) (*entities.Recommendation, error) {
	query := `SELECT id FROM track ORDER BY likes DESC LIMIT 8 OFFSET $1;`
	
	rows, err := r.db.Query(ctx, query, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query rows in ReadNTracksByLikes: %w", err)
	}
	defer rows.Close()

	var trackIds [8]int
	count := 0

	for rows.Next() {
		err = rows.Scan(&trackIds[count])
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows in ReadNTracksByLikes: %w", err)
		}
		count++
	}

	rec := entities.NewRecommendation(trackIds[0], trackIds[1], trackIds[2], trackIds[3], trackIds[4], trackIds[5], trackIds[6], trackIds[7])

	return rec, nil
}
