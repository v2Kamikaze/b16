package auth

import (
	"errors"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrTokenMissing = errors.New("token missing")
var ErrNoPrincipalFound = errors.New("no principal found")
var ErrForbidden = errors.New("forbidden")
