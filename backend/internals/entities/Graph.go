package entities

import "time"

type Graph struct {
	Id          int       `json:"id"`
	TotalTracks int       `json:"totalTracks"`
	CreatedAt   time.Time `json:"createdAt"`
}
