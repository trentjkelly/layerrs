package service

import (
	"context"
	"github.com/trentjkelly/layerr/internals/repository"
	"github.com/trentjkelly/layerr/internals/entities"
)

type RecommendationsService struct {
	trackDbRepo *repository.TrackDatabaseRepository
}

func NewRecommendationsService(trackDbRepo *repository.TrackDatabaseRepository) *RecommendationsService {
	recService := new(RecommendationsService)
	recService.trackDbRepo = trackDbRepo
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
// func (s *RecommendationsService) ArtistLikesAlgorithm(offset int) error {

// }
