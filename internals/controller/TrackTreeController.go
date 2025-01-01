package controller

import (
	"github.com/trentjkelly/layerr/internals/service"
)

type TrackTreeController struct {
	trackTreeService *service.TrackTreeService
}

// Constructor for new TrackTreeController
func NewTrackTreeController(trackTreeService *service.TrackTreeService) *TrackTreeController {
	trackTreeController := new(TrackTreeController)
	trackTreeController.trackTreeService = trackTreeService
	return trackTreeController
}

