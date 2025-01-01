package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerr/internals/config"
	"github.com/trentjkelly/layerr/internals/entities"
)

type ArtistDatabaseRepository struct {
	db	*pgxpool.Pool
}

// Constructor for ArtistDatabaseRepository
func NewArtistDatabaseRepository() *ArtistDatabaseRepository {
	artistDatabaseRepository := new(ArtistDatabaseRepository)
	artistDatabaseRepository.db = config.CreatePSQLPoolConnection()
	return artistDatabaseRepository
}

// Closes the database pool connection
func (r *ArtistDatabaseRepository) CloseDB() {
	r.db.Close()
}

// Adds an Artist to the database, using only the non-optional arguments. 
func (r *ArtistDatabaseRepository) CreateArtist(artist *entities.Artist) {
	query := `INSERT INTO artist (name, username, email) VALUES ($1, $2, $3);`

	row := r.db.QueryRow(context.Background(), query, artist.Name, artist.Username, artist.Email)

	if row != nil {
		fmt.Println(row)
	}
}

// Gets an artist from the database based on their id
func (r *ArtistDatabaseRepository) ReadArtistById(artistId int) {
	query := `SELECT * FROM artist WHERE id=$1;`

	err := r.db.QueryRow(context.Background(), query, artistId)

	if err != nil {
		fmt.Println(err)
	}
}

// Updates the information of the artist
func (r *ArtistDatabaseRepository) UpdateArtist(artist *entities.Artist) {
	query := `UPDATE artist SET (name=$2, email=$3, bio=$4, r2_image_key=$5, updated_at=$6) WHERE id=$1;`

	err := r.db.QueryRow(context.Background(), query, artist.Id, artist.Name, artist.Email, artist.Bio, artist.R2ImageKey, )

	if err != nil {
		fmt.Println(err)
	}
}

// Deletes an artist from the database given the artistId
func (r *ArtistDatabaseRepository) DeleteArtist(artistId int) {
	query := `DELETE FROM artist WHERE id=$1;`

	err := r.db.QueryRow(context.Background(), query, artistId)

	if err != nil {
		fmt.Println(err)
	}
}
