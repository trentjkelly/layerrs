package service

import (
	"github.com/trentjkelly/layerr/internals/repository"
)

type LikesService struct {
	likesDatabaseRepository *repository.LikesDatabaseRepository
}

// Constructor for new LikesService
func NewLikesService(likesDatabaseRepository *repository.LikesDatabaseRepository) *LikesService {
	likesService := new(LikesService)
	likesService.likesDatabaseRepository = likesDatabaseRepository
	return likesService
}