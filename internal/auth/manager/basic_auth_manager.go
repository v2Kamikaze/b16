package manager

import (
	"net/http"

	"github.com/v2code/b16/internal/auth"
	"github.com/v2code/b16/internal/domain"
)

type BasicAuthManager struct {
	users map[string]string
}

func NewBasicAuthManager(users map[string]string) domain.AuthManager[*BasicAuthPrincipal] {
	return &BasicAuthManager{users: users}
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

	expectedPassword, exists := m.users[username]
	if !exists || password != expectedPassword {
		return nil, auth.ErrUnauthorized
	}

	return &BasicAuthPrincipal{Username: username, Password: password}, nil
}
