package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/trentjkelly/layerrs/internals/entities"
	"github.com/trentjkelly/layerrs/internals/service"
)

type PageController struct {
	pageService *service.PageService
}

// Constructor for a new PageController
func NewPageController(pageService *service.PageService) *PageController {
	return &PageController{pageService: pageService}
}

// POST /api/pages — Create a new Page
func (c *PageController) CreatePageHandler(w http.ResponseWriter, r *http.Request) {
	artistId, ok := requireArtistId(w, r)
	if !ok {
		return
	}

	var req entities.CreatePageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Page name is required", http.StatusBadRequest)
		return
	}

	page, err := c.pageService.CreatePage(r.Context(), artistId, req)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	writeJSON(w, page, http.StatusCreated)
}

// GET /api/pages/:pageId — Get Page details
func (c *PageController) GetPageHandler(w http.ResponseWriter, r *http.Request) {
	pageId, err := strconv.Atoi(chi.URLParam(r, "pageId"))
	if err != nil {
		http.Error(w, "Invalid page id", http.StatusBadRequest)
		return
	}

	viewerArtistId := optionalArtistId(r)

	page, err := c.pageService.GetPage(r.Context(), pageId, viewerArtistId)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}
	if page == nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	writeJSON(w, page, http.StatusOK)
}

// PATCH /api/pages/:pageId — Update Page name/description
func (c *PageController) UpdatePageHandler(w http.ResponseWriter, r *http.Request) {
	artistId, ok := requireArtistId(w, r)
	if !ok {
		return
	}

	pageId, err := strconv.Atoi(chi.URLParam(r, "pageId"))
	if err != nil {
		http.Error(w, "Invalid page id", http.StatusBadRequest)
		return
	}

	var req entities.UpdatePageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Page name is required", http.StatusBadRequest)
		return
	}

	page, err := c.pageService.UpdatePage(r.Context(), pageId, artistId, req)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	writeJSON(w, page, http.StatusOK)
}

// DELETE /api/pages/:pageId — Delete Page
func (c *PageController) DeletePageHandler(w http.ResponseWriter, r *http.Request) {
	artistId, ok := requireArtistId(w, r)
	if !ok {
		return
	}

	pageId, err := strconv.Atoi(chi.URLParam(r, "pageId"))
	if err != nil {
		http.Error(w, "Invalid page id", http.StatusBadRequest)
		return
	}

	err = c.pageService.DeletePage(r.Context(), pageId, artistId)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /api/artists/:artistId/pages — Get all Pages owned by an artist
func (c *PageController) GetArtistPagesHandler(w http.ResponseWriter, r *http.Request) {
	artistId, err := strconv.Atoi(chi.URLParam(r, "artistId"))
	if err != nil {
		http.Error(w, "Invalid artist id", http.StatusBadRequest)
		return
	}

	pages, err := c.pageService.GetEditorPages(r.Context(), artistId)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	if pages == nil {
		pages = []entities.Page{}
	}

	writeJSON(w, pages, http.StatusOK)
}

// POST /api/pages/:pageId/follow — Follow a Page
func (c *PageController) FollowPageHandler(w http.ResponseWriter, r *http.Request) {
	artistId, ok := requireArtistId(w, r)
	if !ok {
		return
	}

	pageId, err := strconv.Atoi(chi.URLParam(r, "pageId"))
	if err != nil {
		http.Error(w, "Invalid page id", http.StatusBadRequest)
		return
	}

	err = c.pageService.FollowPage(r.Context(), pageId, artistId)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/pages/:pageId/follow — Unfollow a Page
func (c *PageController) UnfollowPageHandler(w http.ResponseWriter, r *http.Request) {
	artistId, ok := requireArtistId(w, r)
	if !ok {
		return
	}

	pageId, err := strconv.Atoi(chi.URLParam(r, "pageId"))
	if err != nil {
		http.Error(w, "Invalid page id", http.StatusBadRequest)
		return
	}

	err = c.pageService.UnfollowPage(r.Context(), pageId, artistId)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// POST /api/pages/:pageId/tracks — Add a track to a Page (editor only)
func (c *PageController) AddTrackToPageHandler(w http.ResponseWriter, r *http.Request) {
	artistId, ok := requireArtistId(w, r)
	if !ok {
		return
	}

	pageId, err := strconv.Atoi(chi.URLParam(r, "pageId"))
	if err != nil {
		http.Error(w, "Invalid page id", http.StatusBadRequest)
		return
	}

	var req entities.AddTrackToPageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.TrackId == 0 {
		http.Error(w, "trackId is required", http.StatusBadRequest)
		return
	}

	pageTrack, err := c.pageService.AddTrackToPage(r.Context(), pageId, req.TrackId, artistId)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	writeJSON(w, pageTrack, http.StatusCreated)
}

// DELETE /api/pages/:pageId/tracks/:trackId — Remove a track from a Page (editor only)
func (c *PageController) RemoveTrackFromPageHandler(w http.ResponseWriter, r *http.Request) {
	artistId, ok := requireArtistId(w, r)
	if !ok {
		return
	}

	pageId, err := strconv.Atoi(chi.URLParam(r, "pageId"))
	if err != nil {
		http.Error(w, "Invalid page id", http.StatusBadRequest)
		return
	}

	trackId, err := strconv.Atoi(chi.URLParam(r, "trackId"))
	if err != nil {
		http.Error(w, "Invalid track id", http.StatusBadRequest)
		return
	}

	err = c.pageService.RemoveTrackFromPage(r.Context(), pageId, trackId, artistId)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /api/pages/:pageId/feed — Get Page feed (tracks)
func (c *PageController) GetPageFeedHandler(w http.ResponseWriter, r *http.Request) {
	pageId, err := strconv.Atoi(chi.URLParam(r, "pageId"))
	if err != nil {
		http.Error(w, "Invalid page id", http.StatusBadRequest)
		return
	}

	pagination := parsePagination(r)

	pageTracks, err := c.pageService.GetPageFeed(r.Context(), pageId, pagination)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	if pageTracks == nil {
		pageTracks = []entities.PageTrack{}
	}

	writeJSON(w, pageTracks, http.StatusOK)
}

// GET /api/pages/:pageId/submissions — Get submission feed (editor only)
func (c *PageController) GetSubmissionsHandler(w http.ResponseWriter, r *http.Request) {
	artistId, ok := requireArtistId(w, r)
	if !ok {
		return
	}

	pageId, err := strconv.Atoi(chi.URLParam(r, "pageId"))
	if err != nil {
		http.Error(w, "Invalid page id", http.StatusBadRequest)
		return
	}

	pagination := parsePagination(r)

	submissions, err := c.pageService.GetSubmissions(r.Context(), pageId, artistId, pagination)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	if submissions == nil {
		submissions = []entities.PageSubmission{}
	}

	writeJSON(w, submissions, http.StatusOK)
}

// POST /api/pages/:pageId/submissions — Submit a track (followers only)
func (c *PageController) CreateSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	artistId, ok := requireArtistId(w, r)
	if !ok {
		return
	}

	pageId, err := strconv.Atoi(chi.URLParam(r, "pageId"))
	if err != nil {
		http.Error(w, "Invalid page id", http.StatusBadRequest)
		return
	}

	var req entities.SubmitTrackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.TrackId == 0 {
		http.Error(w, "trackId is required", http.StatusBadRequest)
		return
	}

	submission, err := c.pageService.SubmitTrack(r.Context(), pageId, req.TrackId, artistId, req.Note)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	writeJSON(w, submission, http.StatusCreated)
}

// POST /api/pages/:pageId/submissions/:submissionId/approve — Approve a submission (editor only)
func (c *PageController) ApproveSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	artistId, ok := requireArtistId(w, r)
	if !ok {
		return
	}

	submissionId, err := strconv.Atoi(chi.URLParam(r, "submissionId"))
	if err != nil {
		http.Error(w, "Invalid submission id", http.StatusBadRequest)
		return
	}

	var req entities.ApproveSubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	pageTrack, err := c.pageService.ApproveSubmission(r.Context(), submissionId, artistId, req.Note)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	writeJSON(w, pageTrack, http.StatusOK)
}

// GET /api/tracks/:trackId/pages — Get all Pages for a track, ranked
func (c *PageController) GetTrackPagesHandler(w http.ResponseWriter, r *http.Request) {
	trackId, err := strconv.Atoi(chi.URLParam(r, "trackId"))
	if err != nil {
		http.Error(w, "Invalid track id", http.StatusBadRequest)
		return
	}

	response, err := c.pageService.GetTrackPages(r.Context(), trackId)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	writeJSON(w, response, http.StatusOK)
}

// GET /api/feed/following — Get tracks from followed Pages
func (c *PageController) GetFollowingFeedHandler(w http.ResponseWriter, r *http.Request) {
	artistId, ok := requireArtistId(w, r)
	if !ok {
		return
	}

	pagination := parsePagination(r)

	pageTracks, err := c.pageService.GetFollowingFeed(r.Context(), artistId, pagination)
	if err != nil {
		c.handleServiceError(w, err)
		return
	}

	if pageTracks == nil {
		pageTracks = []entities.PageTrack{}
	}

	writeJSON(w, pageTracks, http.StatusOK)
}

// --- Helpers ---

func (c *PageController) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, entities.ErrPageNotFound), errors.Is(err, entities.ErrSubmissionNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, entities.ErrForbidden), errors.Is(err, entities.ErrMustFollowPageToSubmit):
		http.Error(w, err.Error(), http.StatusForbidden)
	case errors.Is(err, entities.ErrDuplicateSubmission), errors.Is(err, entities.ErrConflict):
		http.Error(w, err.Error(), http.StatusConflict)
	default:
		log.Printf("[ERROR] PageController: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func requireArtistId(w http.ResponseWriter, r *http.Request) (int, bool) {
	artistIdFloat, ok := r.Context().Value(entities.ArtistIdKey).(float64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return 0, false
	}
	artistId := int(artistIdFloat)
	if artistId == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return 0, false
	}
	return artistId, true
}

func optionalArtistId(r *http.Request) int {
	artistIdFloat, ok := r.Context().Value(entities.ArtistIdKey).(float64)
	if !ok {
		return 0
	}
	return int(artistIdFloat)
}

func parsePagination(r *http.Request) entities.Pagination {
	limit := 20
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	return entities.Pagination{Limit: limit, Offset: offset}
}

func writeJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("[ERROR] PageController: failed to encode response: %s", err)
	}
}
