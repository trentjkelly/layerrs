package service

import (
	"context"
	"github.com/trentjkelly/layerrs/internals/repository/database"
	"github.com/trentjkelly/layerrs/internals/entities"
)

type RecommendationsService struct {
	trackDbRepo *databaseRepository.TrackDatabaseRepository
	likesDbRepo *databaseRepository.LikesDatabaseRepository
}

func NewRecommendationsService(trackDbRepo *databaseRepository.TrackDatabaseRepository, likesDbRepo *databaseRepository.LikesDatabaseRepository) *RecommendationsService {
	recService := new(RecommendationsService)
	recService.trackDbRepo = trackDbRepo
	recService.likesDbRepo = likesDbRepo
	return recService
}

// Gets the tracks that are the most liked on the entire site all time
func (s *RecommendationsService) MostLikedAlgorithm(ctx context.Context) (*entities.Recommendation, error) {
	rec, err := s.trackDbRepo.ReadNTracksByLikes(ctx, 0)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

// Gets the most recent liked tracks for an artist -- algorithm shown for an individual artist's library page
func (s *RecommendationsService) ArtistLikesAlgorithm(ctx context.Context, artistId int, offset int) ([25]int, error) {
	likesArr, err := s.likesDbRepo.Read25LikesByArtistId(ctx, artistId, offset)
	if err != nil {
		return likesArr, err
	}

	return likesArr, nil
}
