package databaseRepository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerrs/internals/entities"
)

type GenreDatabaseRepository struct {
	db	*pgxpool.Pool
}

// Constructor for a new GenreDatabaseRepository
func NewGenreDatabaseRepository(db *pgxpool.Pool) *GenreDatabaseRepository {
	genreDatabaseRepository := new(GenreDatabaseRepository)
	genreDatabaseRepository.db = db
	return genreDatabaseRepository
}

// Closes the database connection
func (r *GenreDatabaseRepository) CloseDB() {
	r.db.Close()
}

// Creates a genre in the database
func (r *GenreDatabaseRepository) CreateGenre(ctx context.Context, genre *entities.Genre) error {
	query := `INSERT INTO genre (name) VALUES ($1) RETURNING id;`
	row := r.db.QueryRow(ctx, query, genre.Name)
	err := row.Scan(&genre.Id)
	if err != nil {
		return fmt.Errorf("could not create genre in the DB: %w", err)
	}

	return nil
}

// TODO: Not completed
// Creates a mod for the genre in the database
func (r *GenreDatabaseRepository) CreateGenreMod(ctx context.Context, genreMod *entities.GenreMod) error {
	query := `INSERT INTO genre_mods (genre_id, artist_id, is_founder) VALUES ($1, $2, $3) RETURNING id;`
	// TODO: Input is likely something besides 'mod', determine when creating struct
	row := r.db.QueryRow(ctx, query, genreMod.GenreId, genreMod.ArtistId, genreMod.IsFounder)

	err := row.Scan(&genreMod.Id)
	if err != nil {
		return fmt.Errorf("could not create genre mod in the database: %w", err)
	}

	return nil
}

// Adds a track to a genre in the database
func (r *GenreDatabaseRepository) CreateGenreTracks(ctx context.Context, genreTrack *entities.GenreTrack) error {
	query := `INSERT INTO genre_tracks (genre_id, track_id) VALUES ($1, $2) RETURNING id;`
	row := r.db.QueryRow(ctx, query, genreTrack.GenreId, genreTrack.TrackId)
	err := row.Scan(&genreTrack.Id)
	if err != nil {
		return fmt.Errorf("could not create genre tracks in the database: %w", err)
	}

	return nil
}
