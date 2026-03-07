package entities

type Recommendation struct {
	Id           int     `json:"id"`
	Description  string  `json:"description"`
	ArtistId     int     `json:"artistId"`
	ArtistName   string  `json:"artistName"`
	Likes        int     `json:"likes"`
	Layerrs      int     `json:"layerrs"`
	Duration     float64 `json:"duration"`
	WaveformData []int   `json:"waveformData"`
	IsLiked      bool    `json:"isLiked"`
}
