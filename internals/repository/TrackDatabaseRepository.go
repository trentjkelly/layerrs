package repository

import (
	"database/sql"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerr/internals/config"
	"github.com/trentjkelly/layerr/internals/entities"
)

type TrackDatabaseRepository struct {
	db	*pgxpool.Pool
}

// Constructor for TrackDatabaseRepository
func NewTrackDatabaseRepository() *TrackDatabaseRepository {
	trackDatabaseRepository := new(TrackDatabaseRepository)
	trackDatabaseRepository.db = config.CreatePSQLPoolConnection()
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
		return err
	}

	return nil
}

// Gets a Track from the database based on their id
func (r *TrackDatabaseRepository) ReadTrackById(ctx context.Context, track *entities.Track) error {
	query := `SELECT * FROM track WHERE id=$1;`
	row := r.db.QueryRow(ctx, query, track.Id)
	
	// Potential NULL Values
	var r2TrackKey sql.NullString
	var r2CoverKey sql.NullString

	err := row.Scan(&track.Id, &track.Name, &track.ArtistId, &r2TrackKey, &r2CoverKey, &track.CreatedAt, &track.Plays)

	if err != nil {
		return err
	}

	// Potential NULL values converted to empty strings
	if r2TrackKey.Valid {
		track.R2TrackKey = r2TrackKey.String
	} else {
		track.R2TrackKey = ""
	}

	if r2CoverKey.Valid {
		track.R2CoverKey = r2CoverKey.String
	} else {
		track.R2CoverKey = ""
	}

	return nil
}

// Increases the number of plays on a Track when a user plays a song
func (r *TrackDatabaseRepository) IncrementPlays(ctx context.Context, track *entities.Track) error {
	query := `UPDATE track SET plays = plays + 1 WHERE id=$1 RETURNING plays;`
	row := r.db.QueryRow(ctx, query, track.Id)
	err := row.Scan(&track.Plays)
	
	if err != nil {
		return err
	}

	return nil
}

// Updates the information for a Track in the database
func (r *TrackDatabaseRepository) UpdateTrack(ctx context.Context, track *entities.Track) error {
	query := `UPDATE track SET name=$2, r2_track_key=$3, r2_cover_key=$4 WHERE id=$1 RETURNING name;`
	row := r.db.QueryRow(ctx, query, track.Id, track.Name, track.R2TrackKey, track.R2CoverKey)
	err := row.Scan(&track.Name)

	if err != nil {
		return err
	}

	return nil
}

// Deletes a Track from the database given the trackId
func (r *TrackDatabaseRepository) DeleteTrack(ctx context.Context, track *entities.Track) error {
	query := `DELETE FROM track WHERE id=$1 RETURNING id;`
	row := r.db.QueryRow(ctx, query, track.Id)
	err := row.Scan(&track.Id)

	if err != nil {
		return err
	}

	return nil
}
