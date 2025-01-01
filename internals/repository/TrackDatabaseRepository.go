package repository

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerr/internals/config"
	"context"
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

func (r *TrackDatabaseRepository) CreateTrack(trackName string, artistId string) error {
	query := `INSERT INTO track (name, artist_id) VALUES ('@trackName', '@artistId')`
	
	args := pgx.NamedArgs{
		"trackName": trackName,
		"artistId": artistId,
	}

	_, err := r.db.Exec(context.Background(), query, args)

	if err != nil {
		return err
	}

	return nil
}

// func (r *TrackDatabaseRepository) ReadTrackById() error {}

// func (r *TrackDatabaseRepository) UpdateTrack() error {}

// func (r *TrackDatabaseRepository) DeleteTrack() error {}
