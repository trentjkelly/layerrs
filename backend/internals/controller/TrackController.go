package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

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
	trackDescription := strings.TrimSpace(r.FormValue("name"))
	projectDescription := strings.TrimSpace(r.FormValue("projectDescription"))
	if trackDescription == "" {
		trackDescription = strings.TrimSpace(r.FormValue("description"))
	}
	if utf8.RuneCountInString(projectDescription) > 2000 {
		http.Error(w, "Description must be at most 2000 characters", http.StatusBadRequest)
		return
	}
	artistIdFloat := r.Context().Value(entities.ArtistIdKey).(float64)
	layerrsIdStr := r.FormValue("layerrIDs") // Optional - could have no layerrs to credit
	sourceTrackIDsStr := r.FormValue("sourceTrackIDs")
	stemHeaders := r.MultipartForm.File["stemFiles"]

	if trackDescription == "" {
		log.Println("[ERROR] TrackHandlerPost: track description is required")
		http.Error(w, "Track description is required", http.StatusBadRequest)
		return
	}

	if utf8.RuneCountInString(trackDescription) > 100 || utf8.RuneCountInString(trackDescription) < 3 {
		log.Println("[ERROR] TrackHandlerPost: track description must be between 3 and 100 characters")
		http.Error(w, "Track description must be between 3 and 100 characters", http.StatusBadRequest)
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
	var sourceTrackIDs []int
	if sourceTrackIDsStr != "" {
		if err := json.Unmarshal([]byte(sourceTrackIDsStr), &sourceTrackIDs); err != nil {
			http.Error(w, "Invalid source track ID array", http.StatusBadRequest)
			return
		}
	}
	if len(layerrsIdArr) > 1 {
		http.Error(w, "An add-on can target only one project", http.StatusBadRequest)
		return
	}
	if len(stemHeaders)+len(sourceTrackIDs) > 5 {
		http.Error(w, "A project can include at most 5 stems", http.StatusBadRequest)
		return
	}

	var stemNames []string
	if raw := r.FormValue("stemNames"); raw != "" {
		if err := json.Unmarshal([]byte(raw), &stemNames); err != nil || len(stemNames) != len(stemHeaders)+len(sourceTrackIDs) {
			http.Error(w, "Invalid stem names", http.StatusBadRequest)
			return
		}
	} else {
		for _, h := range stemHeaders {
			stemNames = append(stemNames, h.Filename)
		}
		for range sourceTrackIDs {
			stemNames = append(stemNames, "")
		}
	}
	for i, name := range stemNames {
		stemNames[i] = strings.TrimSpace(name)
		if utf8.RuneCountInString(stemNames[i]) > 255 || (r.FormValue("stemNames") != "" && stemNames[i] == "") {
			http.Error(w, "Stem names must be 1–255 characters", http.StatusBadRequest)
			return
		}
	}
	for _, stemHeader := range stemHeaders {
		if !supportedAudio(stemHeader.Filename) {
			http.Error(w, "Stem files must be WAV or FLAC", http.StatusBadRequest)
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
	if !supportedAudio(audioHeader.Filename) {
		http.Error(w, "Audio file must be WAV or FLAC", http.StatusBadRequest)
		return
	}

	// Passing to Service layer
	var createdID int
	err = c.trackService.AddAndUploadTrack(r.Context(), audioFile, audioHeader, stemHeaders, sourceTrackIDs, trackDescription, projectDescription, stemNames, artistIdInt, layerrsIdArr, &createdID)
	if err != nil {
		log.Printf("[ERROR] TrackHandlerPost: %s", err)
		http.Error(w, "Failed to create track", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": createdID})
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

// POST request -- records a play for a given track id (POST /track/{id}/play)
func (c *TrackController) TrackPlayHandlerPost(w http.ResponseWriter, r *http.Request) {
	// Get trackId from request URL
	trackIdStr := chi.URLParam(r, "id")
	trackId, err := strconv.Atoi(trackIdStr)
	if err != nil {
		log.Println("[ERROR] TrackPlayHandlerPost: ", err)
		http.Error(w, "Invalid track id", http.StatusBadRequest)
		return
	}

	// Get artistId from context
	artistIdFloat := r.Context().Value(entities.ArtistIdKey).(float64)
	artistIdInt := int(artistIdFloat)

	// Record the play
	err = c.trackService.RecordPlay(r.Context(), trackId, artistIdInt)
	if err != nil {
		log.Println("[ERROR] TrackPlayHandlerPost: ", err)
		http.Error(w, "Failed to record play", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
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

// GET request -- returns parent track IDs for a batch of track IDs (GET /tracks/uses?trackIds=1,2,3)
func (c *TrackController) TrackUsesHandlerGet(w http.ResponseWriter, r *http.Request) {
	trackIdsParam := r.URL.Query().Get("trackIds")
	if trackIdsParam == "" {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[int][]int{})
		return
	}

	trackIdStrs := strings.Split(trackIdsParam, ",")
	trackIds := make([]int, 0, len(trackIdStrs))
	for _, s := range trackIdStrs {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		trackId, err := strconv.Atoi(s)
		if err != nil {
			log.Println("[ERROR] TrackUsesHandlerGet: invalid track id:", err)
			http.Error(w, "Invalid trackIds parameter", http.StatusBadRequest)
			return
		}
		trackIds = append(trackIds, trackId)
	}

	if len(trackIds) == 0 {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[int][]int{})
		return
	}

	parentsByTrack, err := c.trackService.GetTrackParentsBulk(r.Context(), trackIds)
	if err != nil {
		log.Println("[ERROR] TrackUsesHandlerGet: ", err)
		http.Error(w, "Failed to get track uses", http.StatusInternalServerError)
		return
	}

	if parentsByTrack == nil {
		parentsByTrack = make(map[int][]int)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(parentsByTrack); err != nil {
		log.Println("[ERROR] TrackUsesHandlerGet: ", err)
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

func supportedAudio(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".wav" || ext == ".flac"
}

func (c *TrackController) ProjectDetailHandlerGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}
	project, err := c.trackService.GetProjectDetail(r.Context(), id)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

func (c *TrackController) LegacyStemAudioHandlerGet(w http.ResponseWriter, r *http.Request) {
	projectID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || projectID <= 0 {
		http.Error(w, "Invalid project", http.StatusBadRequest)
		return
	}
	trackID, err := strconv.Atoi(chi.URLParam(r, "trackId"))
	if err != nil || trackID <= 0 {
		http.Error(w, "Invalid stem", http.StatusBadRequest)
		return
	}
	url, err := c.trackService.LegacyStemURL(r.Context(), projectID, trackID)
	if err != nil {
		http.Error(w, "Stem not found", http.StatusNotFound)
		return
	}
	cmd := exec.CommandContext(r.Context(), "ffmpeg", "-v", "error", "-i", url, "-codec:a", "aac", "-b:a", "256k", "-f", "adts", "pipe:1")
	output, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, "Audio unavailable", http.StatusInternalServerError)
		return
	}
	if err := cmd.Start(); err != nil {
		http.Error(w, "Audio unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "audio/aac")
	_, _ = io.Copy(w, output)
	_ = cmd.Wait()
}
