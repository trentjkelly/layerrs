package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerr/internals/config"
)

type TrackTreeDatabaseRepository struct {
	db	*pgxpool.Pool
}

// Constructor for TrackTreeDatabaseRepository
func NewTrackTreeDatabaseRepository() *TrackTreeDatabaseRepository {
	trackTreeDatabaseRepository := new(TrackTreeDatabaseRepository)
	trackTreeDatabaseRepository.db = config.CreatePSQLPoolConnection()
	return trackTreeDatabaseRepository
}

// Closes the database pool connection
func (r *TrackTreeDatabaseRepository) CloseDB() {
	r.db.Close()
}

// func (r *TrackTreeDatabaseRepository) CreateTrackTree() error {}

// func (r *TrackTreeDatabaseRepository) ReadTrackTreeById() error {}

// func (r *TrackTreeDatabaseRepository) UpdateTrackTree() error {}

// func (r *TrackTreeDatabaseRepository) DeleteTrackTree() error {}
