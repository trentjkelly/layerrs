package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerr/internals/config"
	"github.com/trentjkelly/layerr/internals/entities"
)

type LikesDatabaseRepository struct {
	db	*pgxpool.Pool
}

// Constructor for LikesDatabaseRepository
func NewLikesDatabaseRepository() *LikesDatabaseRepository {
	likesDatabaseRepository := new(LikesDatabaseRepository)
	likesDatabaseRepository.db = config.CreatePSQLPoolConnection()
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
func (r *LikesDatabaseRepository) ReadLikesByArtistId(ctx context.Context, artistId int, offset int) ([25]*entities.Like, error) {
	query := `SELECT * FROM artist_likes_track WHERE artist_id=$1 ORDER BY created_at DESC LIMIT 25 OFFSET $2;`
	rows, err := r.db.Query(ctx, query, artistId, offset)

	if err != nil {
		return [25]*entities.Like{}, err
	}
	defer rows.Close()

	// Construct likes array for user
	var likesArray [25]*entities.Like
	i := 0

	for rows.Next() {
		like := new(entities.Like)
		rows.Scan(&like.Id, &like.ArtistId, &like.TrackId, &like.CreatedAt)
		likesArray[i] = like
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
