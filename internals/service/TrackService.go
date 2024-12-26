package service

import (
	"github.com/trentjkelly/layerr/internals/repository"
	"mime/multipart"
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

func (s *TrackService) AddAndUploadTrack(file multipart.File) error {
	
	// TODO: 
	// 1. Create sql table for track
	// 2. Get the unique ID of the track
	// 3. Hash or hex code the track ID for the track name in R2
	// 4. Upload song to R2

	// s.trackStorageRepo.CreateTrack(file)
}