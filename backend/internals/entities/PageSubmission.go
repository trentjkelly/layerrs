package entities

import "time"

type PageSubmission struct {
	Id          int       `json:"id"`
	PageId      int       `json:"pageId"`
	TrackId     int       `json:"trackId"`
	SubmitterId int       `json:"submitterId"`
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"createdAt"`

	// Populated when returning submission feed items
	Track     *Track  `json:"track,omitempty"`
	Submitter *Artist `json:"submitter,omitempty"`
}
