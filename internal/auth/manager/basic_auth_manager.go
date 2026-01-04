package manager

import (
	"net/http"

	"github.com/v2code/b16/internal/auth"
	"github.com/v2code/b16/internal/domain"
)

type BasicAuthManager struct {
	username string
	password string
}

func NewBasicAuthManager(username, password string) domain.AuthManager[*BasicAuthPrincipal] {
	return &BasicAuthManager{username: username, password: password}
}

type BasicAuthPrincipal struct {
	Username string
	Password string
}

func (m *BasicAuthPrincipal) Principal() *BasicAuthPrincipal {
	return m
}

func (m *BasicAuthManager) Authenticate(req *http.Request) (domain.Principal[*BasicAuthPrincipal], error) {
	username, password, ok := req.BasicAuth()
	if !ok {
		return nil, auth.ErrUnauthorized
	}

	if username != m.username || password != m.password {
		return nil, auth.ErrUnauthorized
	}

	return &BasicAuthPrincipal{Username: username, Password: password}, nil
}
