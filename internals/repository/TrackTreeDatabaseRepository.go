package repository

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerr/internals/config"

	"context"
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

// Adds a child-parent relationship between two tracks to the track_tree sql table
func (r *TrackTreeDatabaseRepository) CreateTrackTree(rootTrackId string, childTrackId string) error {
	query := `INSERT INTO track_tree (root_id, child_id) VALUES ('@rootId', '@childId')`

	args := pgx.NamedArgs{
		"rootId": rootTrackId,
		"childId": childTrackId,
	}

	_, err := r.db.Exec(context.Background(), query, args)

	if err != nil {
		return err
	}

	return nil
}

// func (r *TrackTreeDatabaseRepository) ReadTrackTreeById() error {}

// func (r *TrackTreeDatabaseRepository) UpdateTrackTree() error {}

// func (r *TrackTreeDatabaseRepository) DeleteTrackTree() error {}
