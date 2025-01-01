package controller

import (
	"github.com/trentjkelly/layerr/internals/service"
)

type LikesController struct {
	likesService *service.LikesService 
}

// Constructor for a new LikesController
func NewLikesController(likesService *service.LikesService) *LikesController {
	likesController := new(LikesController)
	likesController.likesService = likesService
	return likesController
}



