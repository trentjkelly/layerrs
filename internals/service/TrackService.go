package service

import (
	"github.com/trentjkelly/layerr/internals/repository"
)
type TrackService struct {
	trackStorageRepo *repository.TrackStorageRepository
}

// Constructor for a new TrackService
func NewTrackService(trackStorageRepo *repository.TrackStorageRepository) *TrackService {
	trackService := new(TrackService)
	trackService.trackStorageRepo = trackStorageRepo
	return trackService
}