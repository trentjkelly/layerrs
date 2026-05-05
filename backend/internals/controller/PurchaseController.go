package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/trentjkelly/layerrs/internals/service"
)

type PurchaseController struct {
	purchaseService *service.PurchaseService
}

func NewPurchaseController(purchaseService *service.PurchaseService) *PurchaseController {
	return &PurchaseController{purchaseService: purchaseService}
}

func (c *PurchaseController) CheckoutHandlerOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusNoContent)
}

type CheckoutRequest struct {
	TrackId int `json:"trackId"`
}

func (c *PurchaseController) CheckoutHandlerPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	artistId := r.Context().Value("artistId").(int)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "failed to read request body"}`, http.StatusBadRequest)
		return
	}

	var req CheckoutRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.TrackId <= 0 {
		http.Error(w, `{"error": "invalid track_id"}`, http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(r.Context(), "artistId", artistId)
	response, err := c.purchaseService.CreateCheckoutSession(ctx, artistId, req.TrackId)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *PurchaseController) ConfirmHandlerOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusNoContent)
}

func (c *PurchaseController) ConfirmHandlerPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	artistId := r.Context().Value("artistId").(int)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "failed to read request body"}`, http.StatusBadRequest)
		return
	}

	var req service.ConfirmPurchaseRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.PaymentIntentId == "" {
		http.Error(w, `{"error": "payment_intent_id is required"}`, http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(r.Context(), "artistId", artistId)
	response, err := c.purchaseService.ConfirmPurchase(ctx, artistId, &req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *PurchaseController) GetPurchasesHandlerOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusNoContent)
}

func (c *PurchaseController) GetPurchasesHandlerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	artistId := r.Context().Value("artistId").(int)
	ctx := context.WithValue(r.Context(), "artistId", artistId)

	purchases, err := c.purchaseService.GetUserPurchases(ctx, artistId)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(purchases)
}

func (c *PurchaseController) GetDownloadURLHandlerOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusNoContent)
}

func (c *PurchaseController) GetDownloadURLHandlerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	artistId := r.Context().Value("artistId").(int)
	trackIdStr := chi.URLParam(r, "trackId")

	var trackId int
	_, err := fmt.Sscanf(trackIdStr, "%d", &trackId)
	if err != nil || trackId <= 0 {
		http.Error(w, `{"error": "invalid track_id"}`, http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(r.Context(), "artistId", artistId)
	response, err := c.purchaseService.GetDownloadURL(ctx, artistId, trackId)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
