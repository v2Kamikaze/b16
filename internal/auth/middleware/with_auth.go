package middleware

import (
	"net/http"

	"github.com/v2code/b16/internal/domain"
)

func WithAuth[T any](manager domain.AuthManager[T], handler domain.AuthHandler[T]) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		principal, err := manager.Authenticate(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		handler(w, r, principal)
	}
}
