package entities

import (
	"time"
)

type MagicLinkToken struct {
	Id          int       `json:"id"`
	HashedToken string    `json:"hashedToken"`
	ArtistId    int       `json:"artistId"`
	CreatedAt   time.Time `json:"createdAt"`
}
