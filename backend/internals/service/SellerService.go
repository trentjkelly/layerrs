package service

import (
	"context"
	"fmt"
	"encoding/json"

	"github.com/stripe/stripe-go/v85"
	"github.com/stripe/stripe-go/v85/account"
	"github.com/stripe/stripe-go/v85/accountlink"

	"github.com/trentjkelly/layerrs/internals/entities"
	"github.com/trentjkelly/layerrs/internals/repository/database"
)

type SellerService struct {
	artistDatabaseRepo *databaseRepository.ArtistDatabaseRepository
	returnURL string
	refreshURL string
}

// Constructor for a new SellerService
func NewSellerService(artistDatabaseRepo *databaseRepository.ArtistDatabaseRepository, returnURL string, refreshURL string) *SellerService {
	sellerService := new(SellerService)
	sellerService.artistDatabaseRepo = artistDatabaseRepo
	sellerService.returnURL = returnURL
	sellerService.refreshURL = refreshURL
	return sellerService
}

// Onboards a new seller to Stripe, and returns the onboarding link
func (s *SellerService) OnboardSeller(ctx context.Context, artistId int) (string, error) {
	artist := new(entities.Artist)
	artist.Id = artistId

	err := s.artistDatabaseRepo.ReadArtistById(ctx, artist)
	if err != nil {
		return "", fmt.Errorf("failed to get artist from database: %w", err)
	}

	// If user doesn't have a stripe account, create an express account
	stripeAccountId := artist.StripeAccountId
	if stripeAccountId == "" {
		// First time, create the Express account
        params := &stripe.AccountParams{
            Type:  stripe.String(string(stripe.AccountTypeExpress)),
            Email: stripe.String(artist.Email),
            Capabilities: &stripe.AccountCapabilitiesParams{
                CardPayments: &stripe.AccountCapabilitiesCardPaymentsParams{
                    Requested: stripe.Bool(true),
                },
                Transfers: &stripe.AccountCapabilitiesTransfersParams{
                    Requested: stripe.Bool(true),
                },
            },
        }
        acct, err := account.New(params)
		if err != nil {
			return "", fmt.Errorf("failed to create stripe account: %w", err)
		}

        stripeAccountId = acct.ID

		artist.StripeAccountId = acct.ID
		artist.StripeAccountStatus = "pending"
		err = s.artistDatabaseRepo.UpdateArtistWithStripeAccountDetails(ctx, artist)
		if err != nil {
			return "", fmt.Errorf("failed to update artist stripe account information in database: %w", err)
		}
	}

	linkParams := &stripe.AccountLinkParams{
		Account: stripe.String(stripeAccountId),
		ReturnURL: stripe.String(s.returnURL),
		RefreshURL: stripe.String(s.refreshURL),
		Type: stripe.String("account_onboarding"),
	}
    link, err := accountlink.New(linkParams)
	if err != nil {
		return "", fmt.Errorf("failed to create stripe account link: %w", err)
	}

	return link.URL, nil
}

func (s *SellerService) UpdateSeller(ctx context.Context, stripeEvent stripe.Event) error {
	switch stripeEvent.Type {
    case "account.updated":
        var acct stripe.Account
        err := json.Unmarshal(stripeEvent.Data.Raw, &acct) 
		if err != nil {
            return fmt.Errorf("failed to unmarshal stripe event: %w", err)
        }

        if acct.ChargesEnabled && acct.DetailsSubmitted {
			err := s.artistDatabaseRepo.UpdateArtistWithStripeWebhook(ctx, &entities.Artist{
				StripeAccountId: acct.ID,
				StripeAccountStatus: "active",
			})
			if err != nil {
				return fmt.Errorf("failed to update artist with stripe webhook: %w", err)
			}
        } else if acct.Requirements != nil && acct.Requirements.DisabledReason != "" {
			err := s.artistDatabaseRepo.UpdateArtistWithStripeWebhook(ctx, &entities.Artist{
				StripeAccountId: acct.ID,
				StripeAccountStatus: "restricted",
			})
			if err != nil {
				return fmt.Errorf("failed to update artist with stripe webhook: %w", err)
			}
        }
    }

	return nil
}