package entities

import "time"

type Page struct {
	Id          int       `json:"id"`
	EditorId    *int      `json:"editorId,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	// Computed fields populated by repository/service layers
	FollowerCount int     `json:"followerCount,omitempty"`
	IsFollowing   bool    `json:"isFollowing,omitempty"`
	IsEditor      bool    `json:"isEditor,omitempty"`
	EditorName    *string `json:"editorName,omitempty"`
}

type PageWithFollowerCount struct {
	Page
}

type CreatePageRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdatePageRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AddTrackToPageRequest struct {
	TrackId int `json:"trackId"`
}

type SubmitTrackRequest struct {
	TrackId int    `json:"trackId"`
	Note    string `json:"note"`
}

type ApproveSubmissionRequest struct {
	Note string `json:"note"`
}

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type TrackPagesResponse struct {
	TopPage    *PageWithFollowerCount  `json:"topPage"`
	OtherCount int                     `json:"otherCount"`
	Pages      []PageWithFollowerCount `json:"pages"`
}

type TrackPagesBulkResponse struct {
	Pages     []PageWithFollowerCount `json:"pages"`
	PageCount int                     `json:"pageCount"`
}
