package databaseRepository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerrs/internals/entities"
)

type LikesDatabaseRepository struct {
	db	*pgxpool.Pool
}

// Constructor for LikesDatabaseRepository
func NewLikesDatabaseRepository(db *pgxpool.Pool) *LikesDatabaseRepository {
	likesDatabaseRepository := new(LikesDatabaseRepository)
	likesDatabaseRepository.db = db
	return likesDatabaseRepository
}

// Closes the database pool connection
func (r *LikesDatabaseRepository) CloseDB() {
	r.db.Close()
}

// Creates a like for a Track
func (r *LikesDatabaseRepository) CreateLike(ctx context.Context, like *entities.Like) error {
	query := `INSERT INTO artist_likes_track (artist_id, track_id) VALUES ($1, $2) RETURNING id;`
	row := r.db.QueryRow(ctx, query, like.ArtistId, like.TrackId)
	err := row.Scan(&like.Id)

	if err != nil {
		return err
	}

	return nil
}

// Gets 25 likes sorted most recent to least recent offset by a certain number for a given artist
func (r *LikesDatabaseRepository) Read25LikesByArtistId(ctx context.Context, artistId int, offset int) ([25]int, error) {
	query := `SELECT track_id FROM artist_likes_track WHERE artist_id=$1 ORDER BY created_at DESC LIMIT 25 OFFSET $2;`
	rows, err := r.db.Query(ctx, query, artistId, offset)
	
	var likesArray [25]int

	if err != nil {
		return likesArray, err
	}
	defer rows.Close()

	// Construct likes array for user
	i := 0
	for rows.Next() {
		rows.Scan(&likesArray[i])
		i++
	}
		
	return likesArray, nil
}

func (r *LikesDatabaseRepository) ReadLikeByTrackIdArtistId(ctx context.Context, like *entities.Like) error {
	query := `SELECT id FROM artist_likes_track WHERE artist_id=$1 AND track_id=$2;`
	row := r.db.QueryRow(ctx, query, like.ArtistId, like.TrackId)

	err := row.Scan(&like.Id)
	if err != nil {
		return err
	}

	return nil
}


// Gets full track info for all tracks liked by an artist, sorted most recently liked first
func (r *LikesDatabaseRepository) ReadLikedTracksFullByArtistId(ctx context.Context, artistId int) ([]entities.TrackInfo, error) {
	query := `
		SELECT t.id, t.description, t.artist_id, a.username, a.r2_image_key,
		       t.likes, t.plays, t.layerrs, t.duration, w.waveform_data,
		       CASE WHEN alt.artist_id IS NOT NULL THEN true ELSE false END as is_liked
		FROM artist_likes_track alt
		JOIN track t ON alt.track_id = t.id
		JOIN artist a ON t.artist_id = a.id
		LEFT JOIN waveform w ON w.track_id = t.id
		LEFT JOIN artist_likes_track alt2 ON alt2.track_id = t.id AND alt2.artist_id = $1
		WHERE alt.artist_id = $1 AND t.is_valid = true
		ORDER BY alt.created_at DESC;
	`

	rows, err := r.db.Query(ctx, query, artistId)
	if err != nil {
		return nil, fmt.Errorf("failed to query liked tracks from database: %w", err)
	}
	defer rows.Close()

	var tracks []entities.TrackInfo
	for rows.Next() {
		var track entities.TrackInfo
		var waveformData []int
		var r2ImageKey sql.NullString
		err = rows.Scan(
			&track.Id, &track.Description, &track.ArtistId, &track.ArtistName, &r2ImageKey,
			&track.Likes, &track.Plays, &track.Layerrs, &track.Duration, &waveformData, &track.IsLiked,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan liked tracks from database: %w", err)
		}
		if r2ImageKey.Valid {
			track.R2ImageKey = r2ImageKey.String
		}
		if waveformData != nil {
			track.WaveformData = waveformData
		} else {
			track.WaveformData = []int{}
		}
		tracks = append(tracks, track)
	}
	return tracks, nil
}

// Deletes a like from the database based on the like's artistId & trackId
func (r *LikesDatabaseRepository) DeleteLike(ctx context.Context, like *entities.Like) error {
	query := `DELETE FROM artist_likes_track WHERE artist_id=$1 AND track_id=$2 RETURNING id;`
	row := r.db.QueryRow(ctx, query, like.ArtistId, like.TrackId)
	err := row.Scan(&like.Id)

	if err != nil {
		return err
	}

	return nil
}
