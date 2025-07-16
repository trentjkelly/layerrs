package controller

import (
	"net/http"
	"github.com/trentjkelly/layerrs/internals/service"
	"strconv"
	"encoding/json"
	"github.com/go-chi/chi/v5"
)

type ArtistController struct {
	artistService *service.ArtistService
}

// Constructor for a new ArtistController
func NewArtistController(artistService *service.ArtistService) *ArtistController {
	artistController := new(ArtistController)
	artistController.artistService = artistService
	return artistController
}

// OPTIONS request -- for browsers when they test for CORS before PUT request
func (c *ArtistController) ArtistHandlerOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.WriteHeader(http.StatusOK)
}

// POST request -- creates a new artist
func (c *ArtistController) ArtistHandlerPost(w http.ResponseWriter, r *http.Request) {

}

// PUT request -- updates an existing artist
func (c *ArtistController) ArtistHandlerPut(w http.ResponseWriter, r *http.Request) {

}

// GET request -- Sends the artists informaiton to frontend
func (c *ArtistController) ArtistHandlerGet(w http.ResponseWriter, r *http.Request) {
	// Get artistId
	artistStr := chi.URLParam(r, "artistId")
	artistId, err := strconv.Atoi(artistStr)
	if err != nil {
		http.Error(w, "Invalid artist id", http.StatusBadRequest)
		return
	}

	// Get the rest of the artist data
	artist, err := c.artistService.GetArtistData(r.Context(), artistId)
	if err != nil {
		http.Error(w, "Could not get artist data", http.StatusInternalServerError)
		return
	}

	// Send the data
	err = json.NewEncoder(w).Encode(artist)
	if err != nil {
		http.Error(w, "Could not send artist data", http.StatusInternalServerError)
	}
}

// DELETE request -- Deletes an artist's information
func (c *ArtistController) ArtistHandlerDelete(w http.ResponseWriter, r *http.Request) {

}
