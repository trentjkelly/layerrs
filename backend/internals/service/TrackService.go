package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
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
	trackStorageRepo 		*storageRepository.TrackStorageRepository
	coverStorageRepo 		*storageRepository.CoverStorageRepository
	trackDatabaseRepo 		*databaseRepository.TrackDatabaseRepository
	treeDatabaseRepo 		*databaseRepository.TrackTreeDatabaseRepository
	trackConversionRepo 	*computingRepository.TrackConversionRepository
	waveformHeightsRepo 	*computingRepository.WaveformHeightsRepository
	waveformDatabaseRepo 	*databaseRepository.WaveformDatabaseRepository
	layerrsDatabaseRepo 	*databaseRepository.LayerrsDatabaseRepository
	environment				string
}

// Constructor for a new TrackService
func NewTrackService(
	trackStorageRepo 		*storageRepository.TrackStorageRepository, 
	coverStorageRepo 		*storageRepository.CoverStorageRepository, 
	trackDatabaseRepo 		*databaseRepository.TrackDatabaseRepository, 
	treeDatabaseRepo 		*databaseRepository.TrackTreeDatabaseRepository,
	trackConversionRepo 	*computingRepository.TrackConversionRepository,
	waveformHeightsRepo 	*computingRepository.WaveformHeightsRepository,
	waveformDatabaseRepo 	*databaseRepository.WaveformDatabaseRepository,
	layerrsDatabaseRepo 	*databaseRepository.LayerrsDatabaseRepository,
	environment				string,
) *TrackService {
	trackService := new(TrackService)
	trackService.trackStorageRepo = trackStorageRepo
	trackService.coverStorageRepo = coverStorageRepo
	trackService.trackDatabaseRepo = trackDatabaseRepo
	trackService.treeDatabaseRepo = treeDatabaseRepo
	trackService.trackConversionRepo = trackConversionRepo
	trackService.waveformHeightsRepo = waveformHeightsRepo
	trackService.waveformDatabaseRepo = waveformDatabaseRepo
	trackService.layerrsDatabaseRepo = layerrsDatabaseRepo
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

	// Audio file type conversions
	// TODO: Remove filepaths after being done
	flacPath, opusPath, aacPath, flacName, opusName, aacName, err := s.trackConversionRepo.ConvertAllTracks(audio, trackIdStr, audiofileExtension)
	if err != nil {
		return fmt.Errorf("failed to convert audio file to all formats: %w", err)
	}

	// Add all tracks to R2
	err = s.trackStorageRepo.CreateAllTracks(ctx, flacPath, opusPath, aacPath, flacName, opusName, aacName)
	if err != nil {
		return fmt.Errorf("failed to create all tracks in the storage bucket: %w", err)
	}

	duration, err := s.trackConversionRepo.GetAACTrackDuration(aacPath)
	if err != nil {
		return fmt.Errorf("failed to get the OPUS Track duration: %w", err)
	}

	track.AacR2TrackKey = aacName
	track.FlacR2TrackKey = flacName
	track.OpusR2TrackKey = opusName
	track.TrackDuration = duration

	err = s.trackDatabaseRepo.UpdateTrack(ctx, track)
	if err != nil {
		return fmt.Errorf("failed to update the track in the db: %w", err)
	}

	// Waveform generation
	waveformEntity := new(entities.Waveform)
	waveform, err := s.waveformHeightsRepo.CreateWaveform(flacPath)
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

	// Set is valid flag to true, meaning the track can be served
	track.IsValid = true
	err = s.trackDatabaseRepo.UpdateTrack(ctx, track)
	if err != nil {
		return fmt.Errorf("failed to update the track in the db: %w", err)
	}

	return nil
}

// Gets all of the track's info from the database
func (s *TrackService) GetTrackInfo(ctx context.Context, trackId int) (*entities.Track, error) {
	// Initialize new track
	track := new(entities.Track)
	track.Id = trackId

	waveform := new(entities.Waveform)
	waveform.TrackId = trackId

	// Get the track's info from the database
	err := s.trackDatabaseRepo.ReadTrackById(ctx, track)
	if err != nil {
		return nil, fmt.Errorf("failed to read track from database: %w", err)
	}

	err = s.waveformDatabaseRepo.GetWaveform(ctx, waveform)
	if err != nil {
		return nil, fmt.Errorf("failed to read waveform from database: %w", err)
	}

	track.WaveformData = waveform.WaveformData

	if !track.IsValid {
		return nil, fmt.Errorf("track is not valid")
	}

	return track, nil
}

// Streams a track by its track id
func (s *TrackService) GetStreamingSignedTrackURL(ctx context.Context, trackId int) (string, string, error) {
	track := new(entities.Track)
	track.Id = trackId

	err := s.trackDatabaseRepo.ReadTrackById(ctx, track)
	if err != nil {
		return "", "", fmt.Errorf("failed to read track from database: %w", err)
	}

	if !track.IsValid {
		return "", "", fmt.Errorf("track is not valid")
	}

	url, expiresAt, err := s.trackStorageRepo.GetSignedOpusURL(ctx, track.OpusR2TrackKey, 10*time.Minute)
	if err != nil {
		return "", "", fmt.Errorf("failed to get signed url for track: %w", err)
	}

	return url, expiresAt.String(), nil
}

// Streams a track by its track id
func (s *TrackService) GetDownloadSignedTrackURL(ctx context.Context, trackId int, artistId int) (string, string, error) {
	track := new(entities.Track)
	track.Id = trackId

	// Get track from database
	err := s.trackDatabaseRepo.ReadTrackById(ctx, track)
	if err != nil {
		return "", "", fmt.Errorf("failed to read track from database: %w", err)
	}
	
	if !track.IsValid {
		return "", "", fmt.Errorf("track is not valid")
	}

	// Add track to layerrs list for artist
	layerr := new(entities.Layerr)
	layerr.ArtistId = artistId
	layerr.TrackId = track.Id

	err = s.layerrsDatabaseRepo.CreateLayerr(ctx, layerr)
	if err != nil {
		return "", "", fmt.Errorf("failed to create layerr in the db: %w", err)
	}

	url, expiresAt, err := s.trackStorageRepo.GetSignedFlacURL(ctx, track.FlacR2TrackKey, 1*time.Minute)
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

	if !track.IsValid {
		return nil, fmt.Errorf("track is not valid")
	}

	file, err := s.coverStorageRepo.ReadCover(ctx, &track.R2CoverKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read cover art from storage bucket: %w", err)
	}

	return file, nil
}
