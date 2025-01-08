package controller

import (
	"net/http"
	"encoding/json"
	"github.com/trentjkelly/layerr/internals/service"
)

type RecommendationsController struct {
	recService *service.RecommendationsService
}

func NewRecommendationsController(recService *service.RecommendationsService) *RecommendationsController {
	recController := new(RecommendationsController)
	recController.recService = recService
	return recController
}

// Sends a user what tracks to show on their homepage
func (c *RecommendationsController) RecommendationsHandlerHomeGet(w http.ResponseWriter, r *http.Request) {
	rec, err := c.recService.MostLikedAlgorithm(r.Context())

	if err !=  nil {
		http.Error(w, "Unable to get reccomendations", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(rec)

	if err != nil {
		http.Error(w, "Unable to encode recommendations to json", http.StatusInternalServerError)
	}
}

// Sends a user what tracks to show on their likes page
func (c *RecommendationsController) ReccomendationsHandlerLikesGet(w http.ResponseWriter, r *http.Request) {

}
