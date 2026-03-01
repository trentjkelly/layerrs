package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/trentjkelly/layerrs/internals/service"
	"github.com/trentjkelly/layerrs/internals/entities"
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

// PUT request -- updates an existing artist
func (c *ArtistController) ArtistHandlerPut(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10 * 1024 * 1024) // 10MB max
	skipFile := false

	err := r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		log.Printf("[ERROR] ArtistHandlerPut: %s", err)
		http.Error(w, "Request body too large or malformed", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	if len(username) < 3 || len(username) > 30 {
		log.Println("[ERROR] ArtistHandlerPut: ", "Username is invalid")
		http.Error(w, "Username is invalid", http.StatusBadRequest)
		return
	}

	bio := r.FormValue("bio")
	if len(bio) > 300 {
		log.Println("[ERROR] ArtistHandlerPut: ", "Bio is invalid")
		http.Error(w, "Bio is invalid", http.StatusBadRequest)
		return
	}

	portraitFile, portraitHeader, err := r.FormFile("portraitFile")
	if err != nil && err != http.ErrMissingFile {
		log.Printf("[ERROR] ArtistHandlerPut: %s", err)
		http.Error(w, "Portrait file is required", http.StatusBadRequest)
		return
	}

	if err == http.ErrMissingFile {
		skipFile = true
	} else {
		defer portraitFile.Close()

		ext := filepath.Ext(portraitHeader.Filename)
		allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
		if !allowed[ext] {
			log.Printf("[ERROR] ArtistHandlerPut: invalid file type: %s", ext)
			http.Error(w, "Invalid file type", http.StatusBadRequest)
			return
		}
	}

	artistIdFloat, ok := r.Context().Value(entities.ArtistIdKey).(float64)
	if !ok {
		log.Println("[ERROR] ArtistHandlerPut: could not parse artistId from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	artistId := int(artistIdFloat)
	if artistId == 0 {
		log.Println("[ERROR] ArtistHandlerPut: ", "ArtistId is invalid")
		http.Error(w, "ArtistId is invalid", http.StatusBadRequest)
		return
	}

	err = c.artistService.UpdateArtist(r.Context(), username, bio, portraitFile, portraitHeader, skipFile)
	if err != nil {
		log.Printf("[ERROR] ArtistHandlerPut: %s", err)
		http.Error(w, "Failed to update artist", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GET request -- Sends the artists informaiton to frontend
func (c *ArtistController) ArtistHandlerGet(w http.ResponseWriter, r *http.Request) {
	// Get artistId
	artistStr := chi.URLParam(r, "artistId")
	artistId, err := strconv.Atoi(artistStr)
	if err != nil {
		log.Printf("[ERROR] ArtistHandlerGet: %s", err)
		http.Error(w, "Invalid artist id", http.StatusBadRequest)
		return
	}

	// Get the rest of the artist data
	artist, err := c.artistService.GetArtistData(r.Context(), artistId)
	if err != nil {
		log.Printf("[ERROR] ArtistHandlerGet: %s", err)
		http.Error(w, "Could not get artist data", http.StatusInternalServerError)
		return
	}

	// Send the data
	err = json.NewEncoder(w).Encode(artist)
	if err != nil {
		log.Printf("[ERROR] ArtistHandlerGet: %s", err)
		http.Error(w, "Could not send artist data", http.StatusInternalServerError)
	}
}

// DELETE request -- Deletes an artist's information
func (c *ArtistController) ArtistHandlerDelete(w http.ResponseWriter, r *http.Request) {

}
