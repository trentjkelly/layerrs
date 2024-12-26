package repository

import (
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

// func (r *LikesDatabaseRepository) CreateLikes() error {}

// func (r *LikesDatabaseRepository) ReadLikesById() error {}

// func (r *LikesDatabaseRepository) UpdateLikes() error {}

// func (r *LikesDatabaseRepository) DeleteLikes() error {}
