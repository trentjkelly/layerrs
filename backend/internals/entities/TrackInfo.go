package entities

type TrackInfo struct {
	Id                int     `json:"id"`
	Description       string  `json:"description"`
	ArtistId          int     `json:"artistId"`
	ArtistName        string  `json:"artistName"`
	ArtistPortraitUrl string  `json:"artistPortraitUrl"`
	R2ImageKey        string  `json:"-"`
	Likes             int     `json:"likes"`
	Plays             int     `json:"plays"`
	Layerrs           int     `json:"layerrs"`
	Duration          float64 `json:"duration"`
	WaveformData      []int   `json:"waveformData"`
	IsLiked           bool    `json:"isLiked"`
	PriceInCents      int     `json:"priceInCents"`
}
