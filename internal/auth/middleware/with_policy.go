package middleware

import (
	"net/http"

	"github.com/v2code/b16/internal/domain"
)

func WithPolicy[T any](handler domain.AuthHandler[T], policy domain.Policy[T]) domain.AuthHandler[T] {

	return func(w http.ResponseWriter, r *http.Request, credentials domain.Principal[T]) {

		if err := policy.Check(credentials); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		handler(w, r, credentials)
	}
}
