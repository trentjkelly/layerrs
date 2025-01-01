package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerr/internals/config"
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

func (r *LikesDatabaseRepository) CreateLikes(artistId string, trackId string) error {
	query := `INSERT INTO artist_likes_track (artist_id, track_id) VALUES ('@artistId', '@trackId')`

	args := pgx.NamedArgs{
		"artistId": artistId,
		"trackId": trackId,
	}

	_, err := r.db.Exec(context.Background(), query, args)

	if err != nil {
		return err
	}

	return nil
}

// func (r *LikesDatabaseRepository) ReadLikesById() error {}

// func (r *LikesDatabaseRepository) UpdateLikes() error {}

// func (r *LikesDatabaseRepository) DeleteLikes() error {}
