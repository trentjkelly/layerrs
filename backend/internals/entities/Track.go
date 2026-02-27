package entities

import (
	"time"
)

type Track struct {
	Id 				int			`json:"id"`
	Description 	string		`json:"description"`
	ArtistId 		int			`json:"artistId"`
	FlacR2TrackKey 	string		`json:"flacR2TrackKey"`
	OpusR2TrackKey 	string		`json:"opusR2TrackKey"`
	AacR2TrackKey 	string		`json:"aacR2TrackKey"`
	CreatedAt 		time.Time	`json:"createdAt"`
	Plays			int			`json:"plays"`
	Likes			int 		`json:"likes"`
	Layerrs			int			`json:"layerrs"`
	IsValid			bool		`json:"isValid"`
	WaveformData	[]int		`json:"waveformData"`
	TrackDuration 	float64		`json:"trackDuration"`
}

// Constructor for a new track
func NewTrack(description string, artistId int) *Track {
	track := new(Track)
	track.Description = description
	track.ArtistId = artistId
	return track
}