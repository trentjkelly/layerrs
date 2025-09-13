package controller

import (
	"encoding/json"
	"net/http"

	"github.com/trentjkelly/layerrs/internals/entities"
	"github.com/trentjkelly/layerrs/internals/service"
)

type LayerrsController struct {
	layerrsService *service.LayerrsService
}

// Constructor for a new LayerrsController
func NewLayerrsController(layerrsService *service.LayerrsService) *LayerrsController {
	layerrsController := new(LayerrsController)
	layerrsController.layerrsService = layerrsService
	return layerrsController
}

// OPTIONS request for browsers when they test for CORS before PUT request
func (c *LayerrsController) LayerrsHandlerOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.WriteHeader(http.StatusOK)
}

// GET request -- Sends the artist's layerrs to the frontend
func (c *LayerrsController) LayerrsHandlerGet(w http.ResponseWriter, r *http.Request) {
	// Get artistId
	artistIdFloat := r.Context().Value(entities.ArtistIdKey).(float64)
	artistId := int(artistIdFloat)

	// Get the artist's layerrs
	layerrs, err := c.layerrsService.GetArtistLayerrs(r.Context(), artistId)
	if err != nil {
		http.Error(w, "Could not get artist's layerrs", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(layerrs)
	if err != nil {
		http.Error(w, "Could not encode artist's layerrs to json", http.StatusInternalServerError)
	}
}