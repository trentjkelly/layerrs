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
	query := `INSERT INTO artist_likes_track (artist_id, track_id) VALUES ($1, $2) RETURNING id`
	row := r.db.QueryRow(ctx, query, like.ArtistId, like.TrackId)
	err := row.Scan(&like.Id)

	if err != nil {
		return err
	}

	return nil
}

// Gets all database information about a like by the like id
func (r *LikesDatabaseRepository) ReadLikeById(ctx context.Context, like *entities.Like) error {
	query := `SELECT * FROM artist_likes_track WHERE id=$1`
	row := r.db.QueryRow(ctx, query, like.Id)
	err := row.Scan(&like.Id, &like.ArtistId, &like.TrackId, &like.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

// Deletes a like from the database based on the like's id
func (r *LikesDatabaseRepository) DeleteLikes(ctx context.Context, like *entities.Like) error {
	query := `DELETE FROM artist_likes_track WHERE id=$1 RETURNING id;`
	row := r.db.QueryRow(ctx, query, like.Id)
	err := row.Scan(&like.Id)

	if err != nil {
		return err
	}

	return err
}
