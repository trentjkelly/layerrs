package service

import (
	"github.com/trentjkelly/layerr/internals/repository"
	// "github.com/trentjkelly/layerr/internals/entities"
)

type ArtistService struct {
	artistDatabaseRepository *repository.ArtistDatabaseRepository
}

// Constructor for a new ArtistService
func NewArtistService(artistDatabaseRepository *repository.ArtistDatabaseRepository) *ArtistService {
	artistService := new(ArtistService)
	artistService.artistDatabaseRepository = artistDatabaseRepository
	return artistService
}

// Creates a new artist
func (s *ArtistService) CreateNewArtist() {

	// s.artistDatabaseRepository.CreateArtist()
}