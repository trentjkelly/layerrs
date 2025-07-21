package databaseRepository

import (
	"database/sql"
	"context"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerrs/internals/entities"
	// "log"
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
func (r *ArtistDatabaseRepository) CreateArtist(ctx context.Context, username string, name string, email string, password string) (int, error) {
	query := `INSERT INTO artist (username, name, email, password) VALUES ($1, $2, $3, $4) RETURNING id;`
	row := r.db.QueryRow(ctx, query, username, name, email, password)
	
	var artistId int
	err := row.Scan(&artistId)
	if err != nil {
		return 0, err
	}

	return artistId, nil
}

func (r *ArtistDatabaseRepository) GetArtistIdUsernamePassword(ctx context.Context, artist *entities.Artist, email string) (error) {
	query := `SELECT id, username, password FROM artist WHERE email=$1;`
	row := r.db.QueryRow(ctx, query, email)

	err := row.Scan(&artist.Id, &artist.Username, &artist.Password)
	if err != nil {
		return err
	}

	return nil
}

// Gets an Artist from the database based on their id
func (r *ArtistDatabaseRepository) ReadArtistById(ctx context.Context, artist *entities.Artist) error {
	query := `SELECT * FROM artist WHERE id=$1;`
	row := r.db.QueryRow(ctx, query, artist.Id)

	// Potential NULL Values
	var bio sql.NullString
	var r2ImageKey sql.NullString

	err := row.Scan(&artist.Id, &artist.Name, &artist.Username, &artist.Email, &bio, &r2ImageKey, &artist.CreatedAt, &artist.UpdatedAt, &artist.Password)

	if err != nil {
		return err
	}

	// Potential NULL values converted to empty strings
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

// Updates the information of the artist
func (r *ArtistDatabaseRepository) UpdateArtist(ctx context.Context, artist *entities.Artist) error {
	query := `UPDATE artist SET name=$2, email=$3, bio=$4, r2_image_key=$5, updated_at=$6 WHERE id=$1 RETURNING updated_at;`
	row := r.db.QueryRow(ctx, query, artist.Id, artist.Name, artist.Email, artist.Bio, artist.R2ImageKey, time.Now())
	err := row.Scan(&artist.UpdatedAt)
	
	if err != nil {
		return err
	}

	return nil
}

// Deletes an artist from the database given the artistId
func (r *ArtistDatabaseRepository) DeleteArtist(ctx context.Context, artist *entities.Artist) error {
	query := `DELETE FROM artist WHERE id=$1 RETURNING id;`
	row := r.db.QueryRow(ctx, query, artist.Id)
	err := row.Scan(&artist.Id)

	if err != nil {
		return err
	}

	return nil 
}
