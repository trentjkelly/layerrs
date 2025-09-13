package service

import (
	"context"
	"github.com/trentjkelly/layerrs/internals/entities"
	"github.com/trentjkelly/layerrs/internals/repository/database"
)

type LayerrsService struct {
	layerrsDatabaseRepo *databaseRepository.LayerrsDatabaseRepository
}

// Constructor for a new LayerrsService
func NewLayerrsService(layerrsDatabaseRepo *databaseRepository.LayerrsDatabaseRepository) *LayerrsService {
	layerrsService := new(LayerrsService)
	layerrsService.layerrsDatabaseRepo = layerrsDatabaseRepo
	return layerrsService
}

func (s *LayerrsService) GetArtistLayerrs(ctx context.Context, artistId int) ([]*entities.Layerr, error) {
	layerrs, err := s.layerrsDatabaseRepo.ReadLayerrs(ctx, artistId)
	if err != nil {
		return nil, err
	}
	return layerrs, nil
}