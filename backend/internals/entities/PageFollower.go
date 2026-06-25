package entities

import "time"

type PageFollower struct {
	Id         int       `json:"id"`
	PageId     int       `json:"pageId"`
	ArtistId   int       `json:"artistId"`
	FollowedAt time.Time `json:"followedAt"`
}
