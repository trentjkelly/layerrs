package entities

import (
	"time"
)

type Track struct {
	Id 			int
	Name 		string
	ArtistId 	int
	R2TrackKey 	string
	R2CoverKey 	string
	CreatedAt 	time.Time
	Plays		int
}

// Constructor for a new track
func NewTrack(name string, artistId int) *Track {
	track := new(Track)
	track.Name = name
	track.ArtistId = artistId
	return track
}