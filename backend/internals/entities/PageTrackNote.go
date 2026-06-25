package entities

import "time"

type PageTrackNote struct {
	Id          int       `json:"id"`
	PageTrackId int       `json:"pageTrackId"`
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"createdAt"`
}
