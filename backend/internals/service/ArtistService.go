package service

import (
	"context"
	"fmt"

	"github.com/trentjkelly/layerrs/internals/entities"
	"github.com/trentjkelly/layerrs/internals/repository/database"
)

type ArtistService struct {
	artistDatabaseRepository *databaseRepository.ArtistDatabaseRepository
}

// Constructor for a new ArtistService
func NewArtistService(artistDatabaseRepository *databaseRepository.ArtistDatabaseRepository) *ArtistService {
	artistService := new(ArtistService)
	artistService.artistDatabaseRepository = artistDatabaseRepository
	return artistService
}

func (s *ArtistService) GetArtistData(ctx context.Context, artistId int) (*entities.Artist, error) {
	artist := new(entities.Artist)
	artist.Id = artistId

	err := s.artistDatabaseRepository.ReadArtistById(ctx, artist)
	if err != nil {
		return artist, fmt.Errorf("failed to read artist from database: %w", err)
	}

	return artist, nil
}

// Creates a new artist
func (s *ArtistService) CreateNewArtist() {
	// s.artistDatabaseRepository.CreateArtist()
}