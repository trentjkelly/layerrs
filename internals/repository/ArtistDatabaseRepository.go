package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trentjkelly/layerr/internals/config"
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

// func (r *ArtistDatabaseRepository) CreateArtist() error {}

// func (r *ArtistDatabaseRepository) ReadArtistById() error {}

// func (r *ArtistDatabaseRepository) UpdateArtist() error {}

// func (r *ArtistDatabaseRepository) DeleteArtist() error {}
