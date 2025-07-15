package entities

import (
	"time"
)

type Like struct {
	Id 			int			`json:"id"`
	ArtistId	int			`json:"artistId"`
	TrackId 	int			`json:"trackId"`
	CreatedAt	time.Time	`json:"createdAt"`
}

type LikeCheck struct {
	IsLiked	bool	`json:"isLiked"`
}