package service

import (
	"context"
	"log"
	"time"

	"github.com/trentjkelly/layerrs/internals/entities"
	databaseRepository "github.com/trentjkelly/layerrs/internals/repository/database"
	storageRepository "github.com/trentjkelly/layerrs/internals/repository/storage"
)

type LayerrsService struct {
	layerrsDatabaseRepo *databaseRepository.LayerrsDatabaseRepository
	portraitStorageRepo *storageRepository.PortraitStorageRepository
}

// Constructor for a new LayerrsService
func NewLayerrsService(layerrsDatabaseRepo *databaseRepository.LayerrsDatabaseRepository, portraitStorageRepo *storageRepository.PortraitStorageRepository) *LayerrsService {
	layerrsService := new(LayerrsService)
	layerrsService.layerrsDatabaseRepo = layerrsDatabaseRepo
	layerrsService.portraitStorageRepo = portraitStorageRepo
	return layerrsService
}

func (s *LayerrsService) GetArtistLayerrs(ctx context.Context, artistId int) ([]entities.LayerrTrack, error) {
	layerrTracks, err := s.layerrsDatabaseRepo.ReadLayerrsWithTracks(ctx, artistId)
	if err != nil {
		return nil, err
	}

	for i := range layerrTracks {
		if layerrTracks[i].R2ImageKey != "" {
			url, err := s.portraitStorageRepo.GetSignedPortraitURL(ctx, layerrTracks[i].R2ImageKey, 15*time.Minute)
			if err != nil {
				log.Printf("[WARN] GetArtistLayerrs: could not get signed portrait url: %s", err)
			} else {
				layerrTracks[i].ArtistPortraitUrl = url
			}
		}
	}

	return layerrTracks, nil
}
