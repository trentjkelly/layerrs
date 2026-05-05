package config

import (
	"fmt"
	"os"
)

type StripeConfig struct {
	SecretKey string
	ReturnURL string
	RefreshURL string
	WebhookSecret string
}

// Creates a new Stripe configuration
func NewStripeConfig(env string, frontendURL string) (*StripeConfig, error) {
	stripeConfig := new(StripeConfig)
	stripeSecretKey := os.Getenv(fmt.Sprintf("STRIPE_SECRET_KEY_%s", env))
	if stripeSecretKey == "" {
		return nil, fmt.Errorf("could not find the environment variable STRIPE_SECRET_KEY_%s", env)
	}

	stripeConfig.SecretKey = stripeSecretKey
	stripeConfig.ReturnURL = fmt.Sprintf("%s/profile?onboarding=complete", frontendURL)
	stripeConfig.RefreshURL = fmt.Sprintf("%s/profile?onboarding=refresh", frontendURL)
	stripeConfig.WebhookSecret = os.Getenv(fmt.Sprintf("STRIPE_WEBHOOK_SECRET_%s", env))

	return stripeConfig, nil
}