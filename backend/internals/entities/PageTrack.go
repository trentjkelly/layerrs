package entities

import "time"

type PageTrack struct {
	Id            int       `json:"id"`
	PageId        int       `json:"pageId"`
	TrackId       int       `json:"trackId"`
	AddedAt       time.Time `json:"addedAt"`
	RecommenderId *int      `json:"recommenderId,omitempty"`

	// Populated when returning a feed item
	Track *TrackInfo `json:"track,omitempty"`

	// Editor notes attached to this Page track; only shown on the Page feed
	Notes []PageTrackNote `json:"notes,omitempty"`
}
