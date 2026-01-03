package domain

import (
	"net/http"
)

type Policy[T any] interface {
	Check(credentials UserCredentials[T]) error
}

func WithPolicy[T any](handler AuthHandler[T], policy Policy[T]) AuthHandler[T] {

	return func(w http.ResponseWriter, r *http.Request, credentials UserCredentials[T]) {

		if err := policy.Check(credentials); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		handler(w, r, credentials)
	}
}
