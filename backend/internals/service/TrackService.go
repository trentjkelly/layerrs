package service

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"time"
	"os"

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
	portraitStorageRepo 	*storageRepository.PortraitStorageRepository
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
	portraitStorageRepo 	*storageRepository.PortraitStorageRepository,
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
	trackService.portraitStorageRepo = portraitStorageRepo
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
func (s *TrackService) AddAndUploadTrack(ctx context.Context, audio multipart.File, audioHeader *multipart.FileHeader, trackDescription string, artistId int, parentIDs []int, color string) error {
	// Add track metadata to track table (get back ID)
	track := entities.NewTrack(trackDescription, artistId, color)
	err := s.trackDatabaseRepo.CreateTrack(ctx, track)
	if err != nil {
		return err
	}

	// Update track table with track audio name (the id as a string with .mp3)
	audiofileExtension := filepath.Ext(audioHeader.Filename)
	trackIdStr := strconv.Itoa(track.Id)

	// Audio file type conversions
	flacPath, opusPath, aacPath, flacName, opusName, aacName, err := s.trackConversionRepo.ConvertAllTracks(audio, trackIdStr, audiofileExtension)
	if err != nil {
		return fmt.Errorf("failed to convert audio file to all formats: %w", err)
	}
	defer func() {
		os.Remove(flacPath)
		os.Remove(opusPath)
		os.Remove(aacPath)
	}()

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

	// Create a new graph for the track
	graph := new(entities.Graph)
	graph.TotalTracks = 1
	err = s.treeDatabaseRepo.CreateGraph(ctx, graph)
	if err != nil {
		return err
	}

	// Insert the artist's track into a new graph
	trackGraph := new(entities.TrackGraph)
	trackGraph.TrackId = track.Id
	trackGraph.GraphId = graph.Id
	err = s.treeDatabaseRepo.AddTracktoGraph(ctx, trackGraph)
	if err != nil {
		return fmt.Errorf("failed to add track to graph: %w", err)
	}

	// Track has an array of parents, need to add that relationship to the database as well
	if len(parentIDs) > 0 {
		var trackTrees []*entities.TrackTree
		for _, parentId := range parentIDs {
			trackTree := new(entities.TrackTree)
			trackTree.RootId = parentId
			trackTree.ChildId = track.Id
			trackTree.DerivationTag = "layerr" //TODO: change this later
			trackTrees = append(trackTrees, trackTree)
		}

		err = s.treeDatabaseRepo.CreateGraphRelationships(ctx, trackTrees, track)
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

// Gets a single track's full info by its ID
func (s *TrackService) GetTrackRecommendation(ctx context.Context, trackId int, artistId int) (entities.TrackInfo, error) {
	rec, err := s.trackDatabaseRepo.ReadOneTrackById(ctx, trackId, artistId)
	if err != nil {
		return rec, fmt.Errorf("failed to read track info from database: %w", err)
	}

	if rec.R2ImageKey != "" {
		url, err := s.portraitStorageRepo.GetSignedPortraitURL(ctx, rec.R2ImageKey, 15*time.Minute)
		if err != nil {
			log.Printf("[WARN] GetTrackRecommendation: could not get signed portrait url: %s", err)
		} else {
			rec.ArtistPortraitUrl = url
		}
	}

	return rec, nil
}

// Gets full track info for a batch of track IDs
func (s *TrackService) GetTrackInfoBatch(ctx context.Context, trackIds []int, artistId int) ([]entities.TrackInfo, error) {
	recs, err := s.trackDatabaseRepo.ReadTracksByIds(ctx, trackIds, artistId)
	if err != nil {
		return nil, fmt.Errorf("failed to read tracks from database: %w", err)
	}

	for i, rec := range recs {
		if rec.R2ImageKey != "" {
			url, err := s.portraitStorageRepo.GetSignedPortraitURL(ctx, rec.R2ImageKey, 15*time.Minute)
			if err != nil {
				log.Printf("[WARN] GetTrackInfoBatch: could not get signed portrait url: %s", err)
			} else {
				recs[i].ArtistPortraitUrl = url
			}
		}
	}

	return recs, nil
}

// Gets all TrackTree relationships within the same graph as the given trackId
func (s *TrackService) GetTrackGraphRelationships(ctx context.Context, trackId int) ([]*entities.TrackTree, error) {
	trackTrees, err := s.treeDatabaseRepo.GetGraphTrackTrees(ctx, trackId)
	if err != nil {
		return nil, fmt.Errorf("failed to get graph track trees: %w", err)
	}
	return trackTrees, nil
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
