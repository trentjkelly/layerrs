package repository

import (
	"github.com/trentjkelly/layerr/internals/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type TrackStorageRepository struct {
	r2Config	*aws.Config
	r2Client	*s3.Client
}

// Constructor for new TrackRepository
func NewTrackStorageRepository() *TrackStorageRepository {
	trackStorageRepository := new(TrackStorageRepository)
	trackStorageRepository.r2Config = config.CreateR2Config()
	trackStorageRepository.r2Client = config.CreateR2Client(trackStorageRepository.r2Config)
	return trackStorageRepository
}

// Uploads a track to storage
func (r *TrackStorageRepository) CreateTrack() {

}

// Gets a track from storage (to be streamed)
func (r *TrackStorageRepository) ReadTrack() {

}

// Updates the track in storage
func (r *TrackStorageRepository) UpdateTrack() {

}

// Deletes the track from storage
func (r *TrackStorageRepository) DeleteTrack() {

}


