package service

import (
	"github.com/trentjkelly/layerr/internals/repository"
	"github.com/trentjkelly/layerr/internals/entities"
	"mime/multipart"
	"context"
	"path/filepath"
	"strconv"
)
type TrackService struct {
	trackStorageRepo 	*repository.TrackStorageRepository
	coverStorageRepo 	*repository.CoverStorageRepository
	trackDatabaseRepo 	*repository.TrackDatabaseRepository
	treeDatabaseRepo 	*repository.TrackTreeDatabaseRepository
}

// Constructor for a new TrackService
func NewTrackService(trackStorageRepo *repository.TrackStorageRepository) *TrackService {
	trackService := new(TrackService)
	trackService.trackStorageRepo = trackStorageRepo
	return trackService
}

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
