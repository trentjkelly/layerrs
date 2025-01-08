package controller

import (
	"github.com/trentjkelly/layerr/internals/service"
	"net/http"
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

// Adds a like for a given track from the given artist
func (c *LikesController) LikesHandlerPost(w http.ResponseWriter, r *http.Request) {

}

// Retrieves all liked tracks for a given artist
func (c *LikesController) LikesHandlerGet(w http.ResponseWriter, r *http.Request) {
	
	// artistId := 
	// offset := 

	// likes, err := c.likesService.GetArtistLikes(r.Context(), )

	// if err != nil {
	// 	http.Error(w, "Failed to retrieve likes", http.StatusInternalServerError)
	// 	return
	// }

	// // Encode likes in json and send back
	// for i := 0; i < 25; i++ {

	// }

}

// Removes a like for a given track from the given artist
func (c *LikesController) LikesHandlerDelete(w http.ResponseWriter, r *http.Request) {

}


