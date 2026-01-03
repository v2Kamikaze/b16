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

type BasicAuthParams struct {
	Username string
	Password string
}

func NewBasicAuthManager(params BasicAuthParams) domain.AuthManager[*BasicAuthCredentials] {
	return &BasicAuthManager{username: params.Username, password: params.Password}
}

type BasicAuthCredentials struct {
	Username string
	Password string
}

func (m *BasicAuthCredentials) Principal() *BasicAuthCredentials {
	return m
}

func (m *BasicAuthManager) Authenticate(req *http.Request) (domain.UserCredentials[*BasicAuthCredentials], error) {
	username, password, ok := req.BasicAuth()
	if !ok {
		return nil, auth.ErrUnauthorized
	}

	if username != m.username || password != m.password {
		return nil, auth.ErrUnauthorized
	}

	return &BasicAuthCredentials{Username: username, Password: password}, nil
}
