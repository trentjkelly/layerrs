package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/trentjkelly/layerrs/internals/entities"
	"github.com/trentjkelly/layerrs/internals/service"
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
	// Parse form (for trackAudio file)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Printf("[ERROR] TrackHandlerPost: failed to parse form: %s", err)
		http.Error(w, "Failed to parse form"+err.Error(), http.StatusBadRequest)
		return
	}

	// Getting metadata
	trackDescription := r.FormValue("description")
	artistIdFloat := r.Context().Value(entities.ArtistIdKey).(float64)
	layerrsIdStr := r.FormValue("layerrIDs") // Optional - could have no layerrs to credit


	if trackDescription == "" {
		log.Println("[ERROR] TrackHandlerPost: track description is required")
		http.Error(w, "Track description is required", http.StatusBadRequest)
		return
	}

	if len(trackDescription) > 100 || len(trackDescription) < 10 {
		log.Println("[ERROR] TrackHandlerPost: track description must be between 10 and 100 characters")
		http.Error(w, "Track description must be between 10 and 100 characters", http.StatusBadRequest)
		return
	}

	// Converting layerrsIdStr to integers
	var layerrsIdArr []int
	if layerrsIdStr != "" {
		err := json.Unmarshal([]byte(layerrsIdStr), &layerrsIdArr)
		if err != nil {
			log.Printf("[ERROR] TrackHandlerPost: invalid layerr ID array: %s", err)
			http.Error(w, "Invalid layers ID array", http.StatusBadRequest)
			return
		}
	}

	// Converting artistIdFloat to integer
	artistIdInt := int(artistIdFloat)
	if artistIdInt == 0 || artistIdInt == service.NO_PARENT {
		log.Printf("[ERROR] TrackHandlerPost: invalid artist id: %d", artistIdInt)
		http.Error(w, "Invalid artist id", http.StatusBadRequest)
		return
	}

	// Getting audio file
	audioFile, audioHeader, err := r.FormFile("audioFile")
	if err != nil {
		log.Printf("[ERROR] TrackHandlerPost: audio file is required: %s", err)
		http.Error(w, "Audio file is required", http.StatusBadRequest)
		return
	}
	defer audioFile.Close()

	// Validate that the audio file is in WAV or FLAC format
	audioType := audioHeader.Header.Get("Content-Type")
	if audioType != "audio/wav" && audioType != "audio/flac" {
		log.Printf("[ERROR] TrackHandlerPost: invalid audio type: %s", audioType)
		http.Error(w, "Audio file must be in WAV or FLAC format", http.StatusBadRequest)
		return
	}

	// Passing to Service layer
	err = c.trackService.AddAndUploadTrack(r.Context(), audioFile, audioHeader, trackDescription, artistIdInt, layerrsIdArr)
	if err != nil {
		log.Printf("[ERROR] TrackHandlerPost: %s", err)
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
		log.Println("[ERROR] TrackAudioHandlerGet: ", err)
		http.Error(w, "Invalid track id", http.StatusBadRequest)
		return
	}

	// Get audio from storage
	url, expiresAt, err := c.trackService.GetStreamingSignedTrackURL(r.Context(), trackId)
	if err != nil {
		log.Println("[ERROR] TrackAudioHandlerGet: ", err)
		http.Error(w, "Failed to stream track", http.StatusInternalServerError)
		return
	}

	// Encode response
	var buffer bytes.Buffer
	err = json.NewEncoder(&buffer).Encode(map[string]string{
		"url":       url,
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

// GET request -- streams the audio for a given track id (GET /track/{id}/download)
func (c *TrackController) TrackDownloadHandlerGet(w http.ResponseWriter, r *http.Request) {
	// Get trackId from request URL
	trackIdStr := chi.URLParam(r, "id")
	trackId, err := strconv.Atoi(trackIdStr)
	if err != nil {
		log.Println("[ERROR] TrackAudioHandlerGet: ", err)
		http.Error(w, "Invalid track id", http.StatusBadRequest)
		return
	}

	// Get artistId from context
	artistIdFloat := r.Context().Value(entities.ArtistIdKey).(float64)
	artistIdInt := int(artistIdFloat)

	// Get audio from storage
	url, expiresAt, err := c.trackService.GetDownloadSignedTrackURL(r.Context(), trackId, artistIdInt)
	if err != nil {
		log.Println("[ERROR] TrackAudioHandlerGet: ", err)
		http.Error(w, "Failed to stream track", http.StatusInternalServerError)
		return
	}

	// Encode response
	var buffer bytes.Buffer
	err = json.NewEncoder(&buffer).Encode(map[string]string{
		"url":       url,
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

// POST request -- returns full track info for a list of track IDs (POST /track/batch)
func (c *TrackController) TrackBatchHandlerPost(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TrackIds []int `json:"track_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println("[ERROR] TrackBatchHandlerPost: ", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Optional JWT — extract artistId for isLiked, default 0 for anon
	artistId := 0
	headerString := r.Header.Get("Authorization")
	if headerString != "" {
		parts := strings.Split(headerString, " ")
		if len(parts) == 2 {
			token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("AUTH_SECRET_KEY")), nil
			})
			if err == nil && token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if sub, ok := claims["sub"].(float64); ok {
						artistId = int(sub)
					}
				}
			}
		}
	}

	recs, err := c.trackService.GetTrackInfoBatch(r.Context(), body.TrackIds, artistId)
	if err != nil {
		log.Println("[ERROR] TrackBatchHandlerPost: ", err)
		http.Error(w, "Failed to get track info", http.StatusInternalServerError)
		return
	}

	if recs == nil {
		recs = []entities.TrackInfo{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(recs); err != nil {
		log.Println("[ERROR] TrackBatchHandlerPost: ", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GET request -- returns recommendation-shaped data for a single track (GET /track/{id}/recommendation)
func (c *TrackController) TrackRecommendationHandlerGet(w http.ResponseWriter, r *http.Request) {
	trackIdStr := chi.URLParam(r, "id")
	trackId, err := strconv.Atoi(trackIdStr)
	if err != nil {
		log.Println("[ERROR] TrackRecommendationHandlerGet: ", err)
		http.Error(w, "Invalid track id", http.StatusBadRequest)
		return
	}

	// Optional JWT — extract artistId for isLiked, default 0 for anon
	artistId := 0
	headerString := r.Header.Get("Authorization")
	if headerString != "" {
		parts := strings.Split(headerString, " ")
		if len(parts) == 2 {
			token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("AUTH_SECRET_KEY")), nil
			})
			if err == nil && token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if sub, ok := claims["sub"].(float64); ok {
						artistId = int(sub)
					}
				}
			}
		}
	}

	rec, err := c.trackService.GetTrackRecommendation(r.Context(), trackId, artistId)
	if err != nil {
		log.Println("[ERROR] TrackRecommendationHandlerGet: ", err)
		http.Error(w, "Failed to get track recommendation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rec)
	if err != nil {
		log.Println("[ERROR] TrackRecommendationHandlerGet: ", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GET request -- returns all TrackTree relationships within the graph a track belongs to (GET /track/{id}/graph)
func (c *TrackController) TrackGraphHandlerGet(w http.ResponseWriter, r *http.Request) {
	trackIdStr := chi.URLParam(r, "id")
	trackId, err := strconv.Atoi(trackIdStr)
	if err != nil {
		log.Println("[ERROR] TrackGraphHandlerGet: ", err)
		http.Error(w, "Invalid track id", http.StatusBadRequest)
		return
	}

	trackTrees, err := c.trackService.GetTrackGraphRelationships(r.Context(), trackId)
	if err != nil {
		log.Println("[ERROR] TrackGraphHandlerGet: ", err)
		http.Error(w, "Failed to get track graph relationships", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(trackTrees)
	if err != nil {
		log.Println("[ERROR] TrackGraphHandlerGet: ", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
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
		return
	}
}

type UpdatePriceRequest struct {
	PriceInCents int `json:"priceInCents"`
}

func (c *TrackController) UpdateTrackPriceHandlerOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusNoContent)
}

func (c *TrackController) UpdateTrackPriceHandlerPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	trackIdStr := chi.URLParam(r, "id")
	trackId, err := strconv.Atoi(trackIdStr)
	if err != nil {
		log.Println("[ERROR] UpdateTrackPriceHandlerPut: ", err)
		http.Error(w, "Invalid track id", http.StatusBadRequest)
		return
	}

	artistIdFloat := r.Context().Value(entities.ArtistIdKey).(float64)
	artistIdInt := int(artistIdFloat)

	var req UpdatePriceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("[ERROR] UpdateTrackPriceHandlerPut: ", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = c.trackService.UpdateTrackPrice(r.Context(), trackId, artistIdInt, req.PriceInCents)
	if err != nil {
		log.Println("[ERROR] UpdateTrackPriceHandlerPut: ", err)
		http.Error(w, fmt.Sprintf("Failed to update price: %s", err.Error()), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
