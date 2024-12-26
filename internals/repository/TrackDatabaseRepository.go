package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerr/internals/config"
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

// func (r *TrackDatabaseRepository) CreateTrack() error {}

// func (r *TrackDatabaseRepository) ReadTrackById() error {}

// func (r *TrackDatabaseRepository) UpdateTrack() error {}

// func (r *TrackDatabaseRepository) DeleteTrack() error {}
