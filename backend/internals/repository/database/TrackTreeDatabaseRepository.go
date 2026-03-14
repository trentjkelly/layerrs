package databaseRepository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerrs/internals/entities"
)

type TrackTreeDatabaseRepository struct {
	db	*pgxpool.Pool
}

// Constructor for TrackTreeDatabaseRepository
func NewTrackTreeDatabaseRepository(db *pgxpool.Pool) *TrackTreeDatabaseRepository {
	trackTreeDatabaseRepository := new(TrackTreeDatabaseRepository)
	trackTreeDatabaseRepository.db = db
	return trackTreeDatabaseRepository
}

// Closes the database pool connection
func (r *TrackTreeDatabaseRepository) CloseDB() {
	r.db.Close()
}

func (r *TrackTreeDatabaseRepository) CreateGraphRelationships(ctx context.Context, trackTrees []*entities.TrackTree, track *entities.Track) error {
	// Create database transaction
	options := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}

	tx, err := r.db.BeginTx(ctx, options)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Insert track_tree edges
	parentIds := make([]int, 0, len(trackTrees))
	for _, trackTree := range trackTrees {
		err = r.CreateTrackTree(ctx, tx, trackTree)
		if err != nil {
			return fmt.Errorf("failed to create track tree: %w", err)
		}
		parentIds = append(parentIds, trackTree.RootId)
	}

	// 2. Resolve canonical graph, reassign all tracks, merge totals, delete losers
	err = r.ResolveCanonicalGraph(ctx, tx, parentIds, track.Id)
	if err != nil {
		return fmt.Errorf("failed to resolve canonical graph: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *TrackTreeDatabaseRepository) ResolveCanonicalGraph(ctx context.Context, tx pgx.Tx, parentIds []int, newTrackId int) error {
	// 1. Find canonical graph id (oldest created_at among all parent graphs)
	var canonicalId int
	err := tx.QueryRow(ctx, `
		SELECT g.id FROM graphs g
		JOIN track_graphs tg ON tg.graph_id = g.id
		WHERE tg.track_id = ANY($1)
		ORDER BY g.created_at ASC LIMIT 1
	`, parentIds).Scan(&canonicalId)
	if err != nil {
		return fmt.Errorf("failed to find canonical graph: %w", err)
	}

	// 2. Capture new track's current graph id before reassigning
	var newTrackGraphId int
	err = tx.QueryRow(ctx, `
		SELECT graph_id FROM track_graphs WHERE track_id = $1
	`, newTrackId).Scan(&newTrackGraphId)
	if err != nil {
		return fmt.Errorf("failed to find new track's graph: %w", err)
	}

	// 3. Reassign loser parent track_graphs to canonical
	_, err = tx.Exec(ctx, `
		UPDATE track_graphs
		SET graph_id = $1
		WHERE graph_id IN (
			SELECT g.id FROM graphs g
			JOIN track_graphs tg ON tg.graph_id = g.id
			WHERE tg.track_id = ANY($2) AND g.id != $1
		)
	`, canonicalId, parentIds)
	if err != nil {
		return fmt.Errorf("failed to reassign loser track_graphs: %w", err)
	}

	// 4. Reassign new track to canonical
	_, err = tx.Exec(ctx, `
		UPDATE track_graphs SET graph_id = $1 WHERE track_id = $2
	`, canonicalId, newTrackId)
	if err != nil {
		return fmt.Errorf("failed to reassign new track to canonical graph: %w", err)
	}

	// 5. Merge total_tracks from all loser graphs (parent losers + new track's temp graph) into canonical
	_, err = tx.Exec(ctx, `
		UPDATE graphs
		SET total_tracks = total_tracks + (
			SELECT COALESCE(SUM(total_tracks), 0) FROM graphs
			WHERE id IN (
				SELECT g.id FROM graphs g
				JOIN track_graphs tg ON tg.graph_id = g.id
				WHERE tg.track_id = ANY($2) AND g.id != $1
			)
		) + (
			SELECT total_tracks FROM graphs WHERE id = $3
		)
		WHERE id = $1
	`, canonicalId, parentIds, newTrackGraphId)
	if err != nil {
		return fmt.Errorf("failed to merge total_tracks into canonical graph: %w", err)
	}

	// 6. Delete all loser graphs and the new track's temporary graph
	_, err = tx.Exec(ctx, `
		DELETE FROM graphs WHERE id IN (
			SELECT g.id FROM graphs g
			JOIN track_graphs tg ON tg.graph_id = g.id
			WHERE tg.track_id = ANY($2) AND g.id != $1
			UNION
			SELECT $3
		)
	`, canonicalId, parentIds, newTrackGraphId)
	if err != nil {
		return fmt.Errorf("failed to delete loser graphs: %w", err)
	}

	return nil
}

func (r *TrackTreeDatabaseRepository) CreateGraph(ctx context.Context, graph *entities.Graph) error {
	query := `INSERT INTO graphs (total_tracks) VALUES ($1) RETURNING id;`
	row := r.db.QueryRow(ctx, query, graph.TotalTracks)

	err := row.Scan(&graph.Id)
	if err != nil {
		return fmt.Errorf("failed to create graph: %w", err)
	}

	return nil
}

func (r *TrackTreeDatabaseRepository) AddTracktoGraph(ctx context.Context, trackGraph *entities.TrackGraph) error {
	query := `INSERT INTO track_graphs (track_id, graph_id) VALUES ($1, $2);`

	_, err := r.db.Exec(ctx, query, trackGraph.TrackId, trackGraph.GraphId)
	if err != nil {
		return fmt.Errorf("failed to insert track to graph: %w", err)
	}

	return nil
}

// Creates a child-parent relationship between two tracks to the track_tree sql table
func (r *TrackTreeDatabaseRepository) CreateTrackTree(ctx context.Context, tx pgx.Tx, trackTree *entities.TrackTree) error {
	query := `INSERT INTO track_tree (root_id, child_id, tag) VALUES ($1, $2, $3)`
	_, err := tx.Exec(ctx, query, trackTree.RootId, trackTree.ChildId, trackTree.DerivationTag)
	if err != nil {
		return fmt.Errorf("failed to insert track tree: %w", err)
	}

	return nil
}

// Gets all of the parents of a given track from the database
func (r *TrackTreeDatabaseRepository) GetParents(ctx context.Context, trackTree *entities.TrackTree) error {
	return nil
}

// Gets all of the children of a given track from the database
func (r *TrackTreeDatabaseRepository) GetChildren(ctx context.Context, trackTree *entities.TrackTree) error {
	return nil
}

// Deletes a tree relationship between two tracks from the database
func (r *TrackTreeDatabaseRepository) DeleteTrackTree(ctx context.Context, trackTree *entities.TrackTree) error {
	query := `DELETE FROM track_tree WHERE root_id=$1 AND child_id=$2 RETURNING root_id;`
	row := r.db.QueryRow(ctx, query, trackTree.RootId, trackTree.ChildId)
	
	err := row.Scan(&trackTree.RootId)
	if err != nil {
		return err
	}

	return nil
}
