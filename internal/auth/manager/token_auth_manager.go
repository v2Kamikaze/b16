package manager

import (
	"net/http"
	"strings"

	"github.com/v2code/b16/internal/auth"
	"github.com/v2code/b16/internal/security"
)

type TokenAuthManager struct {
	jwtIssuer security.TokenIssuer[*security.Claims]
}

func NewTokenAuthManager(jwtIssuer security.TokenIssuer[*security.Claims]) *TokenAuthManager {
	return &TokenAuthManager{
		jwtIssuer: jwtIssuer,
	}
}

type TokenPrincipal struct {
	*security.Claims
}

func (p *TokenPrincipal) Principal() *TokenPrincipal {
	return p
}

func (m *TokenAuthManager) Authenticate(req *http.Request) (auth.Principal[*TokenPrincipal], error) {

	authorization := req.Header.Get("Authorization")
	if authorization == "" {
		return nil, auth.ErrUnauthorized
	}

	token := strings.TrimPrefix(authorization, "Bearer ")
	if token == "" {
		return nil, auth.ErrTokenMissing
	}

	claims, err := m.jwtIssuer.Decode(token)
	if err != nil {
		return nil, auth.ErrUnauthorized
	}

	return &TokenPrincipal{
		Claims: claims,
	}, nil
}
