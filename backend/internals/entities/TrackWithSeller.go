package entities

type TrackWithSeller struct {
	TrackId             int    `json:"trackId"`
	ArtistId            int    `json:"artistId"`
	WavR2TrackKey       string `json:"wavR2TrackKey"`
	PriceInCents        int    `json:"priceInCents"`
	StripeAccountId     string `json:"stripeAccountId"`
	StripeAccountStatus string `json:"stripeAccountStatus"`
}
