package controller

import (
	"net/http"
	"strconv"
	"github.com/trentjkelly/layerr/internals/service"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
)

type TrackController struct {
	trackService *service.TrackService
}

// Constructor for a new TrackController
func NewTrackController(trackService *service.TrackService) *TrackController {
	trackController := new(TrackController)
	trackController.trackService = trackService
	return trackController
}

// OPTIONS request for browsers when they test for CORS before PUT request
func (c *TrackController) TrackHandlerOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.WriteHeader(http.StatusNoContent)
}

// POST request -- creating a new track (POST /track)
func (c *TrackController) TrackHandlerPost(w http.ResponseWriter, r *http.Request) {
	log.Println("Request is being handled")

	// Parse form (for trackAudio and coverArt files)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form" + err.Error(), http.StatusBadRequest)
		return
	}

	// Getting metadata
	trackName := r.FormValue("name")
	artistId := r.FormValue("artistId")
	parentIdStr := r.FormValue("parentId") // Optional
	if trackName == "" || artistId == "" {
		http.Error(w, "Track name and artist are required", http.StatusBadRequest)
		return
	}
	parentIdInt := 0
	artistIdInt := 0

	// Converting parentId & artistId to integers
	if parentIdStr != "" {
		parentIdInt, err = strconv.Atoi(parentIdStr)
		if err != nil {
			http.Error(w, "Invalid parent track id", http.StatusBadRequest)
			return
		}
	}
	artistIdInt, err = strconv.Atoi(artistId)
	if err != nil {
		http.Error(w, "Invalid artist id", http.StatusBadRequest)
		return
	}

	// Getting audio file
	audioFile, audioHeader, err := r.FormFile("audioFile")
	if err != nil {
		http.Error(w, "Audio file is required", http.StatusBadRequest)
		return
	}
	defer audioFile.Close()

	// Getting cover art file
	coverArtFile, coverArtHeader, err := r.FormFile("coverArtFile")
	if err != nil {
		http.Error(w, "Cover art file is required", http.StatusBadRequest)
		return
	}
	defer coverArtFile.Close()

	// Passing to Service layer
	err = c.trackService.AddAndUploadTrack(r.Context(), coverArtFile, coverArtHeader, audioFile, audioHeader, trackName, artistIdInt, parentIdInt)
	if err != nil {
		http.Error(w, "Failed to create track", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GET request -- streams the audio for a given track id (GET /track/{id}/audio)
func (c *TrackController) TrackAudioHandlerGet(w http.ResponseWriter, r *http.Request) {
	
	// Get trackId from request URL
	trackIdStr := chi.URLParam(r, "id")
	trackId, err := strconv.Atoi(trackIdStr)
	if err != nil {
		http.Error(w, "Invalid track id", http.StatusBadRequest)
		return
	}

	// Stream audio
	file, err := c.trackService.StreamTrack(r.Context(), trackId)
	if err != nil {
		http.Error(w, "Failed to stream track", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "audio/*")
	w.Header().Set("Content-Disposition", "inline")

	_, err = io.Copy(w, file)

	if err != nil {
		http.Error(w, "Error while streaming audio", http.StatusInternalServerError)
	}
}

func (c *TrackController) TrackCoverHandlerGet(w http.ResponseWriter, r *http.Request) {
	
	// Get trackId from request URL
	trackIdStr := chi.URLParam(r, "id")
	trackId, err := strconv.Atoi(trackIdStr)
	if err != nil {
		http.Error(w, "Invalid track id", http.StatusBadRequest)
		return
	}

	// Send cover back
	file, err := c.trackService.SendCoverArt(r.Context(), trackId)
	if err != nil {
		http.Error(w, "Failed to stream track", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "img/*")

	_, err = io.Copy(w, file)

	if err != nil {
		http.Error(w, "Error while retrieving cover art", http.StatusInternalServerError)
	}
}
