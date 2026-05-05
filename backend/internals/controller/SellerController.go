package controller

import (
	"net/http"
	"encoding/json"
	"log"

	"github.com/trentjkelly/layerrs/internals/service"
	"github.com/trentjkelly/layerrs/internals/entities"
)

type SellerController struct {
	sellerService *service.SellerService
}

// Constructor for a new SellerController
func NewSellerController(sellerService *service.SellerService) *SellerController {
	sellerController := new(SellerController)
	sellerController.sellerService = sellerService
	return sellerController
}

// OPTIONS request for browsers when they test for CORS before PUT request
func (c *SellerController) OnboardSellerHandlerOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.WriteHeader(http.StatusOK)
}

// Onboards a new seller to Stripe
func (c *SellerController) OnboardSellerHandlerPost(w http.ResponseWriter, r *http.Request) {
	artistIdFloat := r.Context().Value(entities.ArtistIdKey).(float64)
	artistId := int(artistIdFloat)

	link, err := c.sellerService.OnboardSeller(r.Context(), artistId)
	if err != nil {
		log.Printf("[ERROR] SellerHandlerPost: failed to onboard seller: %s", err)
		http.Error(w, "Failed to onboard seller", http.StatusInternalServerError)
		return
	}

	log.Printf("[DEBUG] Onboarding link: %s", link)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(link)
	if err != nil {
		log.Printf("[ERROR] SellerHandlerPost: %s", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}