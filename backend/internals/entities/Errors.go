package entities

import (
	"errors"
)

var ErrInvalidToken = errors.New("invalid token")
var ErrUsernameTaken = errors.New("username already taken")