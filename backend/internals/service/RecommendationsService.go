package service

import (
	"context"
	"github.com/trentjkelly/layerrs/internals/repository"
	"github.com/trentjkelly/layerrs/internals/entities"
)

type RecommendationsService struct {
	trackDbRepo *repository.TrackDatabaseRepository
	likesDbRepo *repository.LikesDatabaseRepository
}

func NewRecommendationsService(trackDbRepo *repository.TrackDatabaseRepository, likesDbRepo *repository.LikesDatabaseRepository) *RecommendationsService {
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
