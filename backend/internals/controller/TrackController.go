package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"bytes"

	"github.com/go-chi/chi/v5"
	"github.com/trentjkelly/layerrs/internals/service"
	"github.com/trentjkelly/layerrs/internals/entities"
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
	// Parse form (for trackAudio and coverArt files)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form" + err.Error(), http.StatusBadRequest)
		return
	}

	// Getting metadata
	trackName := r.FormValue("name")
	artistIdFloat := r.Context().Value(entities.ArtistIdKey).(float64)
	parentIdStr := r.FormValue("parentId") // Optional
	if trackName == "" {
		http.Error(w, "Track name is required", http.StatusBadRequest)
		return
	}

	parentIdInt := 0

	// Converting parentId & artistId to integers
	if parentIdStr != "" {
		parentIdInt, err = strconv.Atoi(parentIdStr)
		if err != nil {
			http.Error(w, "Invalid parent track id", http.StatusBadRequest)
			return
		}
	}

	artistIdInt := int(artistIdFloat)
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

	// Validate that the audio file is in WAV or FLAC format
	audioType := audioHeader.Header.Get("Content-Type")
	if audioType != "audio/wav" && audioType != "audio/flac" {
		http.Error(w, "Audio file must be in WAV or FLAC format", http.StatusBadRequest)
		return
	}

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
		log.Println("TrackHandlerPost: ", err)
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
		log.Println("[ERROR] TrackAudioHandlerGet: ", err)
		http.Error(w, "Invalid track id", http.StatusBadRequest)
		return
	}

	// Get audio from storage
	url, expiresAt, err := c.trackService.GetSignedTrackURL(r.Context(), trackId)
	if err != nil {
		log.Println("[ERROR] TrackAudioHandlerGet: ", err)
		http.Error(w, "Failed to stream track", http.StatusInternalServerError)
		return
	}
	
	// Encode response
	var buffer bytes.Buffer
	err = json.NewEncoder(&buffer).Encode(map[string]string{
		"url": url,
		"expiresAt": expiresAt,
	})
	if err != nil {
		log.Println("[ERROR] TrackAudioHandlerGet: ", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Set headers and send response
	w.Header().Set("Content-Type", "application/json")
	w.Write(buffer.Bytes())
}

func (c *TrackController) TrackCoverHandlerGet(w http.ResponseWriter, r *http.Request) {
	
	// Get trackId from request URL
	trackIdStr := chi.URLParam(r, "id")
	trackId, err := strconv.Atoi(trackIdStr)
	if err != nil {
		http.Error(w, "Invalid track id", http.StatusBadRequest)
		return
	}

	// Get cover from storage
	file, err := c.trackService.StreamCoverArt(r.Context(), trackId)
	if err != nil {
		http.Error(w, "Failed to stream track", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Send cover to frontend
	w.Header().Set("Content-Type", "img/*")
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error while retrieving cover art", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *TrackController) TrackerDataHandlerGet(w http.ResponseWriter, r *http.Request) {
		
	// Get trackId from request URL
	trackIdStr := chi.URLParam(r, "id")
	trackId, err := strconv.Atoi(trackIdStr)
	if err != nil {
		log.Println("[ERROR] TrackerDataHandlerGet: ", err)
		http.Error(w, "Invalid track id", http.StatusBadRequest)
		return
	}

	// Get track data from database
	track, err := c.trackService.GetTrackInfo(r.Context(), trackId)
	if err != nil {
		log.Println("[ERROR] TrackerDataHandlerGet: ", err)
		http.Error(w, "Error while getting track data", http.StatusInternalServerError)
		return
	}

	// Encode track and send json 
	err = json.NewEncoder(w).Encode(track)
	if err != nil {
		log.Println("[ERROR] TrackerDataHandlerGet: ", err)
		http.Error(w, "Failed at encoding json", http.StatusInternalServerError)
	}
}
