package service

import (
	"context"
	"io"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"github.com/trentjkelly/layerr/internals/entities"
	"github.com/trentjkelly/layerr/internals/repository"
)
type TrackService struct {
	trackStorageRepo 	*repository.TrackStorageRepository
	coverStorageRepo 	*repository.CoverStorageRepository
	trackDatabaseRepo 	*repository.TrackDatabaseRepository
	treeDatabaseRepo 	*repository.TrackTreeDatabaseRepository
}

// Constructor for a new TrackService
func NewTrackService(trackStorageRepo *repository.TrackStorageRepository, coverStorageRepo *repository.CoverStorageRepository, trackDatabaseRepo *repository.TrackDatabaseRepository, treeDatabaseRepo *repository.TrackTreeDatabaseRepository) *TrackService {
	trackService := new(TrackService)
	trackService.trackStorageRepo = trackStorageRepo
	trackService.coverStorageRepo = coverStorageRepo
	trackService.trackDatabaseRepo = trackDatabaseRepo
	trackService.treeDatabaseRepo = treeDatabaseRepo
	return trackService
}

// Adds all files and data for a new track -- called by TrackController for a POST request
func (s *TrackService) AddAndUploadTrack(ctx context.Context, coverArt multipart.File, coverHeader *multipart.FileHeader, audio multipart.File, audioHeader *multipart.FileHeader, trackName string, artistId int, parentId int) error {

	// Add track metadata to track table (get back ID)
	track := entities.NewTrack(trackName, artistId)
	err := s.trackDatabaseRepo.CreateTrack(ctx, track)

	if err != nil {
		return err
	}

	// Update track table with new cover art and track audio name (the id as a string with .mp3)
	coverFileExtension := filepath.Ext(coverHeader.Filename)
	audiofileExtension := filepath.Ext(audioHeader.Filename)
	trackIdStr := strconv.Itoa(track.Id)

	track.R2CoverKey = trackIdStr + coverFileExtension
	track.R2TrackKey = trackIdStr + audiofileExtension

	err = s.trackDatabaseRepo.UpdateTrack(ctx, track)

	if err != nil {
		return err
	}

	// Add track to track-audio bucket in R2
	err = s.trackStorageRepo.CreateTrack(ctx, audio, &track.R2TrackKey)

	if err != nil {
		return err
	}

	// Add cover art to cover-art bucket in R2
	err = s.coverStorageRepo.CreateCover(ctx, coverArt, &track.R2CoverKey)

	if err != nil {
		return err
	}

	// Track has a parent, need to add that relationship as well
	if parentId != 0 {
		trackTree := new(entities.TrackTree)
		trackTree.RootId = parentId
		trackTree.ChildId = track.Id

		err = s.treeDatabaseRepo.CreateTrackTree(ctx, trackTree)

		if err != nil {
			return err
		}
	}

	// Successful upload
	return nil
}

// Gets all of the track's info from the database
func (s *TrackService) GetTrackInfo(ctx context.Context, trackId int) (*entities.Track, error) {
	
	// Initialize new track
	track := new(entities.Track)
	track.Id = trackId

	// Get the track's info from the database
	err := s.trackDatabaseRepo.ReadTrackById(ctx, track)
	if err != nil {
		return nil, err
	}

	return track, nil
}

// Streams a track by its track id
func (s *TrackService) StreamTrack(ctx context.Context, trackId int) (io.ReadCloser, error) {

	// Get the R2 storage Key
	track := new(entities.Track)
	track.Id = trackId

	err := s.trackDatabaseRepo.ReadTrackById(ctx, track)
	if err != nil {
		return nil, err
	}

	// Stream the track back to the frontend
	file, err := s.trackStorageRepo.ReadTrack(ctx, &track.R2TrackKey)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// Sends a cover back by trackId
func (s *TrackService) StreamCoverArt(ctx context.Context, trackId int) (io.ReadCloser, error) {

	// Get the R2 storage Key
	track := new(entities.Track)
	track.Id = trackId

	err := s.trackDatabaseRepo.ReadTrackById(ctx, track)

	if err != nil {
		return nil, err
	}

	// Stream the track back to the frontend
	file, err := s.coverStorageRepo.ReadCover(ctx, &track.R2CoverKey)

	if err != nil {
		return nil, err
	}

	return file, nil
}
