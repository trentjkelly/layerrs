package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerr/internals/config"
	"github.com/trentjkelly/layerr/internals/entities"
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

// Creates a child-parent relationship between two tracks to the track_tree sql table
func (r *TrackTreeDatabaseRepository) CreateTrackTree(trackTree *entities.TrackTree) error {
	query := `INSERT INTO track_tree (root_id, child_id) VALUES ($1, $2) RETURNING root_id;`
	row := r.db.QueryRow(context.Background(), query, trackTree.RootId, trackTree.ChildId)
	err := row.Scan(&trackTree.RootId)

	if err != nil {
		return err
	}

	return nil
}

// Gets all of the parents of a given track from the database
func (r *TrackTreeDatabaseRepository) GetParents(trackTree *entities.TrackTree) error {
	return nil
}

// Gets all of the children of a given track from the database
func (r *TrackTreeDatabaseRepository) GetChildren(trackTree *entities.TrackTree) error {
	return nil
}

// Deletes a tree relationship between two tracks from the database
func (r *TrackTreeDatabaseRepository) DeleteTrackTree(trackTree *entities.TrackTree) error {
	query := `DELETE FROM track_tree WHERE root_id=$1 AND child_id=$2 RETURNING root_id;`
	row := r.db.QueryRow(context.Background(), query, trackTree.RootId, trackTree.ChildId)
	err := row.Scan(&trackTree.RootId)

	if err != nil {
		return err
	}

	return nil
}
