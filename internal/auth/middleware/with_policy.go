package middleware

import (
	"net/http"

	"github.com/v2code/b16/internal/auth"
)

func WithPolicy[T any](handler auth.AuthHandler[T], policy auth.Policy[T]) auth.AuthHandler[T] {

	return func(w http.ResponseWriter, r *http.Request, principal auth.Principal[T]) {

		if err := policy.Check(principal); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		handler(w, r, principal)
	}
}
