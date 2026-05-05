package entities

import (
	"time"
)

type Track struct {
	Id            int       `json:"id"`
	Description   string    `json:"description"`
	ArtistId      int       `json:"artistId"`
	WavR2TrackKey string    `json:"wavR2TrackKey"`
	AacR2TrackKey string    `json:"aacR2TrackKey"`
	CreatedAt     time.Time `json:"createdAt"`
	Plays         int       `json:"plays"`
	Likes         int       `json:"likes"`
	Layerrs       int       `json:"layerrs"`
	IsValid       bool      `json:"isValid"`
	WaveformData  []int     `json:"waveformData"`
	TrackDuration float64   `json:"trackDuration"`
	Color         string    `json:"color"`
	PriceInCents  int       `json:"priceInCents"`
}

// Constructor for a new track
func NewTrack(description string, artistId int, color string) *Track {
	track := new(Track)
	track.Description = description
	track.ArtistId = artistId
	track.Color = color
	return track
}
