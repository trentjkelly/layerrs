package controller

import (
	"io"
	"log"
	"net/http"

	"github.com/trentjkelly/layerrs/internals/service"
	"github.com/stripe/stripe-go/v85/webhook"
)

type WebhookController struct {
	stripeWebhookSecret string
	sellerService *service.SellerService	
}

// Constructor for a new WebhookController
func NewWebhookController(sellerService *service.SellerService, stripeWebhookSecret string) *WebhookController {
	webhookController := new(WebhookController)
	webhookController.sellerService = sellerService
	webhookController.stripeWebhookSecret = stripeWebhookSecret
	return webhookController
}

// OPTIONS request for browsers when they test for CORS before PUT request
func (c *WebhookController) StripeWebhookHandlerOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.WriteHeader(http.StatusOK)
}

// POST request -- receives webhooks from Stripe
func (c *WebhookController) StripeWebhookHandlerPost(w http.ResponseWriter, r *http.Request) {
    payload, err := io.ReadAll(r.Body)
    if err != nil {
		log.Printf("[ERROR] WebhookHandlerPost: failed to read request body: %s", err)
        http.Error(w, "Bad request", 400)
        return
    }

	sigHeader := r.Header.Get("Stripe-Signature")
    event, err := webhook.ConstructEventWithOptions(payload, sigHeader, c.stripeWebhookSecret, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})
    if err != nil {
		log.Printf("[ERROR] WebhookHandlerPost: failed to construct event: %s", err)
        http.Error(w, "Invalid signature", 400)
        return
    }

	err = c.sellerService.UpdateSeller(r.Context(), event)
	if err != nil {
		log.Printf("[ERROR] WebhookHandlerPost: failed to update seller: %s", err)
		http.Error(w, "Failed to update seller", 400)
		return
	}

	w.WriteHeader(http.StatusOK)
}