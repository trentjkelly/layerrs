package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/trentjkelly/layerrs/internals/entities"
	"github.com/trentjkelly/layerrs/internals/repository/computing"
	"github.com/trentjkelly/layerrs/internals/repository/database"
	"github.com/trentjkelly/layerrs/internals/repository/storage"
)

const (
	NO_PARENT = 0
)
type TrackService struct {
	trackStorageRepo 	*storageRepository.TrackStorageRepository
	coverStorageRepo 	*storageRepository.CoverStorageRepository
	trackDatabaseRepo 	*databaseRepository.TrackDatabaseRepository
	treeDatabaseRepo 	*databaseRepository.TrackTreeDatabaseRepository
	trackConversionRepo *computingRepository.TrackConversionRepository
	waveformComputingRepo 		*computingRepository.WaveformComputingRepository
	waveformDatabaseRepo *databaseRepository.WaveformDatabaseRepository
	environment			string
}

// Constructor for a new TrackService
func NewTrackService(
	trackStorageRepo 		*storageRepository.TrackStorageRepository, 
	coverStorageRepo 		*storageRepository.CoverStorageRepository, 
	trackDatabaseRepo 		*databaseRepository.TrackDatabaseRepository, 
	treeDatabaseRepo 		*databaseRepository.TrackTreeDatabaseRepository,
	trackConversionRepo 	*computingRepository.TrackConversionRepository,
	waveformComputingRepo 	*computingRepository.WaveformComputingRepository,
	waveformDatabaseRepo 	*databaseRepository.WaveformDatabaseRepository,
	environment				string,
) *TrackService {
	trackService := new(TrackService)
	trackService.trackStorageRepo = trackStorageRepo
	trackService.coverStorageRepo = coverStorageRepo
	trackService.trackDatabaseRepo = trackDatabaseRepo
	trackService.treeDatabaseRepo = treeDatabaseRepo
	trackService.trackConversionRepo = trackConversionRepo
	trackService.waveformComputingRepo = waveformComputingRepo
	trackService.waveformDatabaseRepo = waveformDatabaseRepo
	trackService.environment = environment
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
	artistIdStr := strconv.Itoa(artistId)

	// Audio file type conversions
	// TODO: Remove filepaths after being done
	flacPath, opusPath, aacPath, flacName, opusName, aacName, err := s.trackConversionRepo.ConvertAllFormats(audio, audiofileExtension, artistIdStr, trackIdStr)
	if err != nil {
		return fmt.Errorf("failed to convert audio file to all formats: %w", err)
	}

	// Add all tracks to R2
	err = s.trackStorageRepo.CreateAllTracks(ctx, flacPath, opusPath, aacPath)
	if err != nil {
		return fmt.Errorf("failed to create all tracks in the storage bucket: %w", err)
	}

	track.AacR2TrackKey = aacName
	track.FlacR2TrackKey = flacName
	track.OpusR2TrackKey = opusName

	err = s.trackDatabaseRepo.UpdateTrack(ctx, track)
	if err != nil {
		return fmt.Errorf("failed to update the track in the db: %w", err)
	}

	// Add cover art to cover-art bucket in R2
	// err = s.coverStorageRepo.CreateCover(ctx, coverArt, &track.R2CoverKey)
	// if err != nil {
	// 	return fmt.Errorf("failed to create the cover art in the storage bucket: %w", err)
	// }

	// Waveform generation
	audioFile, err := os.Open(flacPath)
	if err != nil {
		return fmt.Errorf("failed to open flac file for creating a waveform: %w", err)
	}
	defer audioFile.Close()

	waveformEntity := new(entities.Waveform)
	waveform, err := s.waveformComputingRepo.CreateWaveform(ctx, audioFile)
	if err != nil {
		return fmt.Errorf("failed to create the waveform for the audio file: %w", err)
	}

	waveformEntity.TrackId = track.Id
	waveformEntity.WaveformData = waveform

	err = s.waveformDatabaseRepo.CreateWaveform(ctx, waveformEntity)
	if err != nil {
		return fmt.Errorf("failed to create the waveform in the db: %w", err)
	}

	// Track has a parent, need to add that relationship as well
	if parentId != NO_PARENT {
		trackTree := new(entities.TrackTree)
		trackTree.RootId = parentId
		trackTree.ChildId = track.Id

		err = s.treeDatabaseRepo.CreateTrackTree(ctx, trackTree)
		if err != nil {
			return fmt.Errorf("failed to create the track tree in the db: %w", err)
		}
	}

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
		return nil, fmt.Errorf("failed to read track from database: %w", err)
	}

	return track, nil
}

// Streams a track by its track id
func (s *TrackService) GetSignedTrackURL(ctx context.Context, trackId int) (string, string, error) {
	track := new(entities.Track)
	track.Id = trackId

	err := s.trackDatabaseRepo.ReadTrackById(ctx, track)
	if err != nil {
		return "", "", fmt.Errorf("failed to read track from database: %w", err)
	}

	url, expiresAt, err := s.trackStorageRepo.GetSignedURL(ctx, &track.FlacR2TrackKey, 10*time.Minute)
	if err != nil {
		return "", "", fmt.Errorf("failed to get signed url for track: %w", err)
	}

	return url, expiresAt.String(), nil
}

// Sends a cover back by trackId
func (s *TrackService) StreamCoverArt(ctx context.Context, trackId int) (io.ReadCloser, error) {
	track := new(entities.Track)
	track.Id = trackId

	err := s.trackDatabaseRepo.ReadTrackById(ctx, track)
	if err != nil {
		return nil, fmt.Errorf("failed to read cover art from storage bucket: %w", err)
	}

	file, err := s.coverStorageRepo.ReadCover(ctx, &track.R2CoverKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read cover art from storage bucket: %w", err)
	}

	return file, nil
}
