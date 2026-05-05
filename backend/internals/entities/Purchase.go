package entities

import (
	"time"
)

type Purchase struct {
	Id                    int       `json:"id"`
	BuyerId               int       `json:"buyerId"`
	TrackId               int       `json:"trackId"`
	StripePaymentIntentId string    `json:"stripePaymentIntentId"`
	AmountCents           int       `json:"amountCents"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}
