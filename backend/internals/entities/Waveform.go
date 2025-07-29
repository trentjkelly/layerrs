package entities

type Waveform struct {
	Id int `json:"id"`
	TrackId int `json:"trackId"`
	WaveformData []int `json:"waveformData"`
}