package service

import (
	"github.com/trentjkelly/layerr/internals/repository"
)

type TrackTreeService struct {
	trackTreeDatabaseRepository *repository.TrackTreeDatabaseRepository
}

// Constructor for a new TrackTreeService
func NewTrackTreeService(trackTreeDatabaseRepository *repository.TrackTreeDatabaseRepository) *TrackTreeService {
	trackTreeService := new(TrackTreeService)
	trackTreeService.trackTreeDatabaseRepository = trackTreeDatabaseRepository
	return trackTreeService
}