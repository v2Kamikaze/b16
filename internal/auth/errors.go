package auth

import (
	"errors"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrTokenMissing = errors.New("token missing")
var ErrNoCredentialsFound = errors.New("no credentials found")
var ErrForbidden = errors.New("forbidden")
