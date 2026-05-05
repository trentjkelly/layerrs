package entities

import (
	"time"
)

type Artist struct {
	Id                  int       `json:"id"`
	Username            string    `json:"username"`
	Email               string    `json:"email"`
	Bio                 string    `json:"bio"`
	R2ImageKey          string    `json:"r2ImageKey"`
	PortraitUrl         string    `json:"portraitUrl,omitempty"`
	CanPost             bool      `json:"canPost"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
	StripeAccountId     string    `json:"stripeAccountId,omitempty"`
	StripeAccountStatus string    `json:"stripeAccountStatus,omitempty"`
}

type ArtistDTO struct {
	Id                  int       `json:"id"`
	Username            string    `json:"username"`
	Email               string    `json:"email"`
	Bio                 string    `json:"bio"`
	PortraitUrl         string    `json:"portraitUrl,omitempty"`
	CanPost             bool      `json:"canPost"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
	ValidStripeSeller 	bool      `json:"isValidStripeSeller"`
}
