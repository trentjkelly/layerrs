package entities

import (
	"time"
)

type Artist struct {
	Id 			int
	Name 		string
	Username 	string
	Email 		string
	Bio 		string
	R2ImageKey 	string
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}

type Like struct {
	Id 			int
	ArtistId	int
	TrackId 	int
	CreatedAt	time.Time
}

type TrackTree struct {
	RootId		int
	ChildId 	int
}