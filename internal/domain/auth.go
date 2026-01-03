package domain

import (
	"net/http"
)

type UserCredentials[T any] interface {
	Principal() T
}

type AuthManager[T any] interface {
	Authenticate(req *http.Request) (UserCredentials[T], error)
}

type AuthHandler[T any] func(w http.ResponseWriter, r *http.Request, credentials UserCredentials[T])

func Auth[T any](manager AuthManager[T], handler AuthHandler[T]) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		cred, err := manager.Authenticate(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		handler(w, r, cred)
	}
}
