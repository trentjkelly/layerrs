package entities

import (
	"time"
)

type Track struct {
	Id 			int			`json:"id"`
	Name 		string		`json:"name"`
	ArtistId 	int			`json:"artistId"`
	R2TrackKey 	string		`json:"r2TrackKey"`
	R2CoverKey 	string		`json:"r2CoverKey"`
	CreatedAt 	time.Time	`json:"createdAt"`
	Plays		int			`json:"plays"`
	Likes		int 		`json:"likes"`
}

// Constructor for a new track
func NewTrack(name string, artistId int) *Track {
	track := new(Track)
	track.Name = name
	track.ArtistId = artistId
	return track
}