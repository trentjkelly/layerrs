package service

import (
	"context"
	"fmt"
	"io"
	"bytes"

	"github.com/trentjkelly/layerrs/internals/entities"
	"github.com/trentjkelly/layerrs/internals/repository/database"
	"github.com/trentjkelly/layerrs/internals/repository/storage"
	"github.com/trentjkelly/layerrs/internals/repository/computing"
	"mime/multipart"
)

type ArtistService struct {
	artistDatabaseRepository  *databaseRepository.ArtistDatabaseRepository
	portraitStorageRepo 	 *storageRepository.PortraitStorageRepository
	portraitConversionRepo 	 *computingRepository.PortraitConversionRepository
}

// Constructor for a new ArtistService
func NewArtistService(artistDatabaseRepository *databaseRepository.ArtistDatabaseRepository, portraitStorageRepository *storageRepository.PortraitStorageRepository, computingRepository *computingRepository.PortraitConversionRepository) *ArtistService {	
	artistService := new(ArtistService)
	artistService.artistDatabaseRepository = artistDatabaseRepository
	artistService.portraitStorageRepo = portraitStorageRepository
	artistService.portraitConversionRepo = computingRepository
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

// Updates the artist's information
func (s *ArtistService) UpdateArtist(ctx context.Context, username string, bio string, portraitFile multipart.File, portraitHeader *multipart.FileHeader, skipFile bool) error {
	var filename string
	
	artistIdFloat, ok := ctx.Value(entities.ArtistIdKey).(float64)
	if !ok {
		return fmt.Errorf("could not parse artistId from context")
	}
	artistId := int(artistIdFloat)
	if artistId == 0 {
		return fmt.Errorf("artistId is invalid")
	}

	if !skipFile {
		filename = fmt.Sprintf("%d.webp", artistId)

		portraitBytes, err := io.ReadAll(portraitFile)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
	
		webpBytes, err := s.portraitConversionRepo.ConvertToWebP(portraitBytes)
		if err != nil {
			return fmt.Errorf("failed to convert file to webp: %w", err)
		}

		body := bytes.NewReader(webpBytes)
		err = s.portraitStorageRepo.CreatePortrait(ctx, body, filename)
		if err != nil {
			return fmt.Errorf("failed to upload portrait to storage: %w", err)
		}
	}

	artist := new(entities.Artist)
	artist.Id = artistId
	artist.Username = username
	artist.Bio = bio
	artist.R2ImageKey = filename
	
	// Update the database with the new data
	var err error
	if skipFile {
		err = s.artistDatabaseRepository.UpdateArtist(ctx, artist)
	} else {
		err = s.artistDatabaseRepository.UpdateArtistWithPortrait(ctx, artist)
	}
	if err != nil {
		return fmt.Errorf("failed to update artist in database: %w", err)
	}

	return nil
}