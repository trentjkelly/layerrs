package entities

import "time"

type Layerr struct {
	Id 			int			`json:"id"`
	ArtistId	int			`json:"artistId"`
	TrackId 	int			`json:"trackId"`
	LastLayerrAt	time.Time	`json:"lastLayerrAt"`
}