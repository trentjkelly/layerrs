package databaseRepository

import (
	"database/sql"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerrs/internals/entities"
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

	err := row.Scan(&track.Id, &track.Name, &track.ArtistId, &r2TrackKey, &r2CoverKey, &track.CreatedAt, &track.Plays, &track.Likes, &track.Layerrs)
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

// Increases the number of likes on a Track when a user likes it
func (r *TrackDatabaseRepository) IncrementLikes(ctx context.Context, track *entities.Track) error {
	query := `UPDATE track SET likes = likes + 1 WHERE id=$1 RETURNING likes;`
	row := r.db.QueryRow(ctx, query, track.Id)
	
	err := row.Scan(&track.Likes)
	if err != nil {
		return err
	}

	return nil
}

// Decreases the number of likes on a Track when a user likes it
func (r *TrackDatabaseRepository) DecrementLikes(ctx context.Context, track *entities.Track) error {
	query := `UPDATE track SET likes = likes - 1 WHERE id=$1 RETURNING likes;`
	row := r.db.QueryRow(ctx, query, track.Id)
	
	err := row.Scan(&track.Likes)
	if err != nil {
		return err
	}

	return nil
}

// Gets the top N tracks by likes -- used for recommendations algorithm
func (r *TrackDatabaseRepository) ReadNTracksByLikes(ctx context.Context, offset int) (*entities.Recommendation, error) {
	query := `SELECT id FROM track ORDER BY likes DESC LIMIT 8 OFFSET $1;`
	
	rows, err := r.db.Query(ctx, query, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trackIds [8]int
	count := 0

	for rows.Next() {
		err = rows.Scan(&trackIds[count])
		if err != nil {
			return nil, err
		}
		count++
	}

	rec := entities.NewRecommendation(trackIds[0], trackIds[1], trackIds[2], trackIds[3], trackIds[4], trackIds[5], trackIds[6], trackIds[7])

	return rec, nil
}
