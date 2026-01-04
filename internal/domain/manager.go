package domain

import (
	"net/http"
)

type Principal[T any] interface {
	Principal() T
}

type AuthManager[T any] interface {
	Authenticate(req *http.Request) (Principal[T], error)
}

type AuthHandler[T any] func(w http.ResponseWriter, r *http.Request, principal Principal[T])
