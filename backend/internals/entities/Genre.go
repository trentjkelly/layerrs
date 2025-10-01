package entities

import "time"

type Genre struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type GenreMod struct {
	Id int `json:"id"`
	GenreId int `json:"genreId"`
	ArtistId int `json:"artistId"`
	IsFounder bool `json:"isFounder"`
	AddedAt time.Time `json:"addedAt"`
}

type GenreTrack struct {
	Id int `json:"id"`
	GenreId int `json:"genreId"`
	TrackId int `json:"trackId"`
	AddedAt time.Time `json:"addedAt"`
}