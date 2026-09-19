package entities

import "time"

type Project struct {
	Id           int       `json:"id"`
	ArtistId     int       `json:"artistId"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"createdAt"`
	Plays        int       `json:"plays"`
	Likes        int       `json:"likes"`
	Layerrs      int       `json:"layerrs"`
	PriceInCents int       `json:"priceInCents"`
	IsValid      bool      `json:"isValid"`
}

type ProjectTrackRole string

const (
	ProjectMasterTrack ProjectTrackRole = "master"
	ProjectStemTrack   ProjectTrackRole = "stem"
)
