package service

import (
	"context"
	"log"
	"time"

	"github.com/trentjkelly/layerrs/internals/entities"
	databaseRepository "github.com/trentjkelly/layerrs/internals/repository/database"
	storageRepository "github.com/trentjkelly/layerrs/internals/repository/storage"
)

type RecommendationsService struct {
	trackDbRepo         *databaseRepository.TrackDatabaseRepository
	likesDbRepo         *databaseRepository.LikesDatabaseRepository
	portraitStorageRepo *storageRepository.PortraitStorageRepository
}

func NewRecommendationsService(trackDbRepo *databaseRepository.TrackDatabaseRepository, likesDbRepo *databaseRepository.LikesDatabaseRepository, portraitStorageRepo *storageRepository.PortraitStorageRepository) *RecommendationsService {
	recService := new(RecommendationsService)
	recService.trackDbRepo = trackDbRepo
	recService.likesDbRepo = likesDbRepo
	recService.portraitStorageRepo = portraitStorageRepo
	return recService
}

// Gets the tracks that are the most liked on the entire site all time
func (s *RecommendationsService) MostLikedAlgorithm(ctx context.Context) ([]entities.TrackInfo, error) {
	recs, err := s.trackDbRepo.ReadNTracksByLikes(ctx, 0)
	if err != nil {
		return nil, err
	}

	for i := range recs {
		if recs[i].R2ImageKey != "" {
			url, err := s.portraitStorageRepo.GetSignedPortraitURL(ctx, recs[i].R2ImageKey, 15*time.Minute)
			if err != nil {
				log.Printf("[WARN] MostLikedAlgorithm: could not get signed portrait url: %s", err)
			} else {
				recs[i].ArtistPortraitUrl = url
			}
		}
	}

	return recs, nil
}

func (s *RecommendationsService) MostRecentAlgorithm(ctx context.Context, artistId int) ([]entities.TrackInfo, error) {
	recs, err := s.trackDbRepo.ReadNTracksByDate(ctx, 0, artistId)
	if err != nil {
		return nil, err
	}

	for i := range recs {
		if recs[i].R2ImageKey != "" {
			url, err := s.portraitStorageRepo.GetSignedPortraitURL(ctx, recs[i].R2ImageKey, 15*time.Minute)
			if err != nil {
				log.Printf("[WARN] MostRecentAlgorithm: could not get signed portrait url: %s", err)
			} else {
				recs[i].ArtistPortraitUrl = url
			}
		}
	}

	return recs, nil
}

// Gets the most recent liked tracks for an artist with full track info
func (s *RecommendationsService) ArtistLikesAlgorithm(ctx context.Context, artistId int) ([]entities.TrackInfo, error) {
	tracks, err := s.likesDbRepo.ReadLikedTracksFullByArtistId(ctx, artistId)
	if err != nil {
		return nil, err
	}

	for i := range tracks {
		if tracks[i].R2ImageKey != "" {
			url, err := s.portraitStorageRepo.GetSignedPortraitURL(ctx, tracks[i].R2ImageKey, 15*time.Minute)
			if err != nil {
				log.Printf("[WARN] ArtistLikesAlgorithm: could not get signed portrait url: %s", err)
			} else {
				tracks[i].ArtistPortraitUrl = url
			}
		}
	}

	return tracks, nil
}
