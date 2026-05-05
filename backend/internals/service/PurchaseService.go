package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/stripe/stripe-go/v85"
	"github.com/stripe/stripe-go/v85/paymentintent"
	"github.com/trentjkelly/layerrs/internals/entities"
	"github.com/trentjkelly/layerrs/internals/repository/database"
	storageRepository "github.com/trentjkelly/layerrs/internals/repository/storage"
)

type PurchaseService struct {
	purchaseRepo     *databaseRepository.PurchaseRepository
	trackStorageRepo *storageRepository.TrackStorageRepository
	returnURL        string
	refreshURL       string
}

func NewPurchaseService(purchaseRepo *databaseRepository.PurchaseRepository, trackStorageRepo *storageRepository.TrackStorageRepository, returnURL string, refreshURL string) *PurchaseService {
	return &PurchaseService{
		purchaseRepo:     purchaseRepo,
		trackStorageRepo: trackStorageRepo,
		returnURL:        returnURL,
		refreshURL:       refreshURL,
	}
}

type CheckoutResponse struct {
	PaymentIntentId string `json:"paymentIntentId"`
	ClientSecret    string `json:"clientSecret"`
	AmountCents     int    `json:"amountCents"`
}

func (s *PurchaseService) CreateCheckoutSession(ctx context.Context, buyerId int, trackId int) (*CheckoutResponse, error) {
	trackWithSeller, err := s.purchaseRepo.ReadTrackWithSeller(ctx, trackId)
	if err != nil {
		return nil, fmt.Errorf("failed to get track: %w", err)
	}
	if trackWithSeller == nil {
		return nil, fmt.Errorf("track not found")
	}

	if trackWithSeller.PriceInCents <= 0 {
		return nil, fmt.Errorf("track is not for sale")
	}

	if trackWithSeller.StripeAccountStatus != "active" || trackWithSeller.StripeAccountId == "" {
		return nil, fmt.Errorf("seller has not completed Stripe setup")
	}

	if trackWithSeller.WavR2TrackKey == "" {
		return nil, fmt.Errorf("track has no WAV file available")
	}

	hasPurchased, err := s.purchaseRepo.HasPurchased(ctx, buyerId, trackId)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing purchase: %w", err)
	}
	if hasPurchased {
		return nil, fmt.Errorf("you have already purchased this track")
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(trackWithSeller.PriceInCents)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		TransferData: &stripe.PaymentIntentTransferDataParams{
			Destination: stripe.String(trackWithSeller.StripeAccountId),
		},
		Metadata: map[string]string{
			"buyer_id":  fmt.Sprintf("%d", buyerId),
			"track_id":  fmt.Sprintf("%d", trackId),
			"seller_id": fmt.Sprintf("%d", trackWithSeller.ArtistId),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	return &CheckoutResponse{
		PaymentIntentId: pi.ID,
		ClientSecret:    pi.ClientSecret,
		AmountCents:     int(pi.Amount),
	}, nil
}

type ConfirmPurchaseRequest struct {
	PaymentIntentId string `json:"paymentIntentId"`
}

type ConfirmPurchaseResponse struct {
	DownloadURL                  string `json:"downloadURL"`
	DownloadURLExpirationMinutes int    `json:"downloadUrlExpirationMinutes"`
}

func (s *PurchaseService) ConfirmPurchase(ctx context.Context, buyerId int, req *ConfirmPurchaseRequest) (*ConfirmPurchaseResponse, error) {
	pi, err := paymentintent.Get(req.PaymentIntentId, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}

	if pi.Status != stripe.PaymentIntentStatusSucceeded {
		return nil, fmt.Errorf("payment not completed, status: %s", pi.Status)
	}

	buyerIdStr, trackIdStr, ok := pi.Metadata["buyer_id"], pi.Metadata["track_id"], true
	if !ok || buyerIdStr != fmt.Sprintf("%d", buyerId) {
		return nil, fmt.Errorf("invalid payment intent metadata")
	}

	trackId := 0
	fmt.Sscanf(trackIdStr, "%d", &trackId)

	existingPurchase, err := s.purchaseRepo.ReadPurchaseByPaymentIntentId(ctx, req.PaymentIntentId)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing purchase: %w", err)
	}
	if existingPurchase != nil {
		return s.generateDownloadURL(ctx, trackId)
	}

	purchase := &entities.Purchase{
		BuyerId:               buyerId,
		TrackId:               trackId,
		StripePaymentIntentId: req.PaymentIntentId,
		AmountCents:           int(pi.Amount),
	}

	err = s.purchaseRepo.CreatePurchase(ctx, purchase)
	if err != nil {
		return nil, fmt.Errorf("failed to record purchase: %w", err)
	}

	return s.generateDownloadURL(ctx, trackId)
}

func (s *PurchaseService) generateDownloadURL(ctx context.Context, trackId int) (*ConfirmPurchaseResponse, error) {
	trackWithSeller, err := s.purchaseRepo.ReadTrackWithSeller(ctx, trackId)
	if err != nil || trackWithSeller == nil {
		return nil, fmt.Errorf("track not found")
	}

	if trackWithSeller.WavR2TrackKey == "" {
		return nil, fmt.Errorf("track has no WAV file")
	}

	expirationMinutes := 15
	downloadURL, _, err := s.trackStorageRepo.GetSignedWavURL(ctx, trackWithSeller.WavR2TrackKey, time.Duration(expirationMinutes)*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to generate download URL: %w", err)
	}

	return &ConfirmPurchaseResponse{
		DownloadURL:                  downloadURL,
		DownloadURLExpirationMinutes: expirationMinutes,
	}, nil
}

type Purchase struct {
	Id               int    `json:"id"`
	TrackId          int    `json:"trackId"`
	TrackDescription string `json:"trackDescription"`
	ArtistName       string `json:"artistName"`
	AmountCents      int    `json:"amountCents"`
	CreatedAt        string `json:"createdAt"`
}

func (s *PurchaseService) GetUserPurchases(ctx context.Context, buyerId int) ([]Purchase, error) {
	purchases, err := s.purchaseRepo.ReadPurchasesByBuyerId(ctx, buyerId)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchases: %w", err)
	}

	var result []Purchase
	for _, p := range purchases {
		result = append(result, Purchase{
			Id:          p.Id,
			TrackId:     p.TrackId,
			AmountCents: p.AmountCents,
			CreatedAt:   p.CreatedAt.Format("2006-01-02"),
		})
	}

	return result, nil
}

func (s *PurchaseService) GetDownloadURL(ctx context.Context, buyerId int, trackId int) (*ConfirmPurchaseResponse, error) {
	hasPurchased, err := s.purchaseRepo.HasPurchased(ctx, buyerId, trackId)
	if err != nil {
		return nil, fmt.Errorf("failed to check purchase: %w", err)
	}
	if !hasPurchased {
		return nil, fmt.Errorf("you have not purchased this track")
	}

	return s.generateDownloadURL(ctx, trackId)
}

func (s *PurchaseService) HandleWebhookEvent(ctx context.Context, stripeEvent stripe.Event) error {
	switch stripeEvent.Type {
	case "payment_intent.succeeded":
		var pi stripe.PaymentIntent
		err := json.Unmarshal(stripeEvent.Data.Raw, &pi)
		if err != nil {
			return fmt.Errorf("failed to unmarshal payment intent: %w", err)
		}

		buyerIdStr, trackIdStr := pi.Metadata["buyer_id"], pi.Metadata["track_id"]
		if buyerIdStr == "" || trackIdStr == "" {
			return nil
		}

		buyerId, trackId := 0, 0
		fmt.Sscanf(buyerIdStr, "%d", &buyerId)
		fmt.Sscanf(trackIdStr, "%d", &trackId)

		existingPurchase, err := s.purchaseRepo.ReadPurchaseByPaymentIntentId(ctx, pi.ID)
		if err != nil {
			return fmt.Errorf("failed to check existing purchase: %w", err)
		}
		if existingPurchase != nil {
			return nil
		}

		purchase := &entities.Purchase{
			BuyerId:               buyerId,
			TrackId:               trackId,
			StripePaymentIntentId: pi.ID,
			AmountCents:           int(pi.Amount),
		}

		return s.purchaseRepo.CreatePurchase(ctx, purchase)
	}
	return nil
}
