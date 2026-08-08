package entities

import (
	"errors"
)

var ErrInvalidToken = errors.New("invalid token")
var ErrUsernameTaken = errors.New("username already taken")

// Page errors
var ErrPageNotFound = errors.New("page not found")
var ErrSubmissionNotFound = errors.New("submission not found")
var ErrForbidden = errors.New("forbidden")
var ErrMustFollowPageToSubmit = errors.New("must follow page to submit")
var ErrDuplicateSubmission = errors.New("duplicate submission")
var ErrConflict = errors.New("conflict")
var ErrPageNameTooLong = errors.New("page name must be 255 characters or less")
