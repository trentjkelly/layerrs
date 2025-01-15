package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"log"

	"github.com/trentjkelly/layerrs/internals/entities"
	"github.com/trentjkelly/layerrs/internals/service"
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

func (c *LikesController) LikesHandlerOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.WriteHeader(http.StatusOK)
}

// Adds a like for a given track from the given artist
func (c *LikesController) LikesHandlerPost(w http.ResponseWriter, r *http.Request) {

	// Get artistId & trackId
	artistIdFloat := r.Context().Value(entities.ArtistIdKey).(float64)
	artistId := int(artistIdFloat)

	trackStr := r.FormValue("trackId")
	trackId, err := strconv.Atoi(trackStr)
	if err != nil {
		http.Error(w, "Invalid track id", http.StatusBadRequest)
		return
	}

	// Add the like
	err = c.likesService.AddLike(r.Context(), artistId, trackId)
	if err != nil {
		http.Error(w, "Could not add like to track", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Returns whether the artist has liked a given track
func (c *LikesController) LikesHandlerGet(w http.ResponseWriter, r *http.Request) {

	// Get artistId & trackId
	artistIdFloat := r.Context().Value(entities.ArtistIdKey).(float64)
	artistId := int(artistIdFloat)
	trackStr := r.URL.Query().Get("trackId")
	trackId, err := strconv.Atoi(trackStr)
	if err != nil {
		http.Error(w, "Invalid trackId", http.StatusBadRequest)
		return
	}

	// Check if like is present (errorr eturned if it's not)
	likeCheck := new(entities.LikeCheck)
	likeCheck.IsLiked = true
	err = c.likesService.CheckLike(r.Context(), artistId, trackId)
	if err != nil {
		likeCheck.IsLiked = false
	}

	// Send json
	err = json.NewEncoder(w).Encode(likeCheck)
	if err != nil {
		http.Error(w, "Could not send back like", http.StatusInternalServerError)
	}
}

// Removes a like for a given track from the given artist
func (c *LikesController) LikesHandlerDelete(w http.ResponseWriter, r *http.Request) {

	// Get artistId & trackId
	artistIdFloat := r.Context().Value(entities.ArtistIdKey).(float64)
	artistId := int(artistIdFloat)

	trackStr := r.URL.Query().Get("trackId")
	trackId, err := strconv.Atoi(trackStr)
	if err != nil {
		http.Error(w, "Invalid trackId", http.StatusBadRequest)
		return
	}
	
	// Delete the like
	err = c.likesService.RemoveLike(r.Context(), artistId, trackId)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not remove a like for the track", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}


