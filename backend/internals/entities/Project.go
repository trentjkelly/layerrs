package entities

import "time"

type Project struct {
	Id           int       `json:"id"`
	ArtistId     int       `json:"artistId"`
	Name         string    `json:"name"`
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

// ProjectDetail describes the public project and its ordered audio assets.
type ProjectDetail struct {
	Id          int            `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	ArtistName  string         `json:"artistName"`
	ParentId    *int           `json:"parentId"`
	ParentName  *string        `json:"parentName"`
	Tracks      []ProjectAudio `json:"tracks"`
}
type ProjectAudio struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	ArtistName string `json:"artistName"`
	URL        string `json:"url"`
	Key        string `json:"-"`
	Legacy     bool   `json:"-"`
}
