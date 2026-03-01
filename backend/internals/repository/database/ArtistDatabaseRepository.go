package databaseRepository

import (
	"database/sql"
	"context"
	"fmt"
	"time"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerrs/internals/entities"
)

type ArtistDatabaseRepository struct {
	db	*pgxpool.Pool
}

// Constructor for ArtistDatabaseRepository
func NewArtistDatabaseRepository(db *pgxpool.Pool) *ArtistDatabaseRepository {
	artistDatabaseRepository := new(ArtistDatabaseRepository)
	artistDatabaseRepository.db = db
	return artistDatabaseRepository
}

// Closes the database pool connection
func (r *ArtistDatabaseRepository) CloseDB() {
	r.db.Close()
}

// Adds an Artist to the database, but only the non-optional fields 
func (r *ArtistDatabaseRepository) CreateArtist(ctx context.Context, email string) (int, error) {
	query := `INSERT INTO artist (email) VALUES ($1) RETURNING id;`
	row := r.db.QueryRow(ctx, query, email)
	
	var artistId int
	err := row.Scan(&artistId)
	if err != nil {
		return 0, fmt.Errorf("failed to create artist in database: %w", err)
	}

	return artistId, nil
}

func (r *ArtistDatabaseRepository) GetArtistByEmail(ctx context.Context, email string) (*entities.Artist, error) {
	query := `SELECT id, username, email FROM artist WHERE email=$1;`
	row := r.db.QueryRow(ctx, query, email)

	var artist entities.Artist
	var username sql.NullString
	
	err := row.Scan(&artist.Id, &username, &artist.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query artist from database: %w", err)
	}
	if username.Valid {
		artist.Username = username.String
	}

	return &artist, nil
}

// Gets an Artist from the database based on their id
func (r *ArtistDatabaseRepository) ReadArtistById(ctx context.Context, artist *entities.Artist) error {
	query := `SELECT id, username, email, bio, r2_image_key, created_at, updated_at FROM artist WHERE id=$1;`
	row := r.db.QueryRow(ctx, query, artist.Id)

	// Potential NULL Values
	var username sql.NullString
	var bio sql.NullString
	var r2ImageKey sql.NullString

	err := row.Scan(&artist.Id, &username, &artist.Email, &bio, &r2ImageKey, &artist.CreatedAt, &artist.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to query artist from database: %w", err)
	}

	// Potential NULL values converted to empty strings
	if username.Valid {
		artist.Username = username.String
	}
	if bio.Valid {
		artist.Bio = bio.String
	} else {
		artist.Bio = ""
	}

	if r2ImageKey.Valid {
		artist.R2ImageKey = r2ImageKey.String
	} else {
		artist.R2ImageKey = ""
	}

	return nil
}

// Updates the information of the artist without changing the portrait key
func (r *ArtistDatabaseRepository) UpdateArtist(ctx context.Context, artist *entities.Artist) error {
	query := `UPDATE artist SET username=$2, bio=$3, updated_at=$4 WHERE id=$1 RETURNING updated_at;`
	row := r.db.QueryRow(ctx, query, artist.Id, artist.Username, artist.Bio, time.Now())
	err := row.Scan(&artist.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update artist in database: %w", err)
	}

	return nil
}

// Updates the information of the artist including the portrait key
func (r *ArtistDatabaseRepository) UpdateArtistWithPortrait(ctx context.Context, artist *entities.Artist) error {
	log.Printf("Updating artist with portrait: %s", artist.Username)
	query := `UPDATE artist SET username=$2, bio=$3, r2_image_key=$4, updated_at=$5 WHERE id=$1 RETURNING updated_at;`
	row := r.db.QueryRow(ctx, query, artist.Id, artist.Username, artist.Bio, artist.R2ImageKey, time.Now())
	err := row.Scan(&artist.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update artist in database: %w", err)
	}

	return nil
}
