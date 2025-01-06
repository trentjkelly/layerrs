package service

import (
	"github.com/trentjkelly/layerr/internals/repository"
	"github.com/trentjkelly/layerr/internals/entities"
	"context"
)

type LikesService struct {
	likesDatabaseRepository *repository.LikesDatabaseRepository
	trackDatabaseRepository *repository.TrackDatabaseRepository
}

// Constructor for new LikesService
func NewLikesService(likesDatabaseRepository *repository.LikesDatabaseRepository, trackDatabaseRepository *repository.TrackDatabaseRepository) *LikesService {
	likesService := new(LikesService)
	likesService.likesDatabaseRepository = likesDatabaseRepository
	likesService.trackDatabaseRepository = trackDatabaseRepository
	return likesService
}

func (s *LikesService) AddLike(ctx context.Context, artistId int, trackId int) error {
	
	// Adds a row to the artist_likes_track table
	like := new(entities.Like)
	like.ArtistId = artistId
	like.TrackId = trackId

	err := s.likesDatabaseRepository.CreateLike(ctx, like)

	if err != nil {
		return err
	}

	// Increments the like counter for a given track (track table)
	track := new(entities.Track)
	track.Id = trackId

	err = s.trackDatabaseRepository.IncrementLikes(ctx, track)

	if err != nil {
		return err
	}

	return nil
}

func (s *LikesService) GetArtistLikes(ctx context.Context, artistId int, offset int) ([25]*entities.Like, error) {
	// Returns all of the tracks an artist has liked, sorted by most to least recently
	likes, err := s.likesDatabaseRepository.ReadLikesByArtistId(ctx, artistId, offset)

	if err != nil {
		return [25]*entities.Like{}, err
	}

	return likes, nil
}

func (s *LikesService) RemoveLike() error {

	// Deletes the row from the artist_likes_track table
	// s.likesDatabaseRepository.RemoveLike()

	// Decrements the like counter for a given track (track table)

	return nil
}