package manager

import (
	"net/http"
	"strings"

	"github.com/v2code/b16/internal/auth"
	"github.com/v2code/b16/internal/domain"
	"github.com/v2code/b16/internal/security"
)

type TokenAuthManager struct {
	jwtIssuer domain.TokenIssuer[*security.Claims]
}

func NewTokenAuthManager(jwtIssuer domain.TokenIssuer[*security.Claims]) *TokenAuthManager {
	return &TokenAuthManager{
		jwtIssuer: jwtIssuer,
	}
}

type TokenCredentials struct {
	*security.Claims
}

func (c *TokenCredentials) Principal() *TokenCredentials {
	return c
}

func (m *TokenAuthManager) Authenticate(req *http.Request) (domain.UserCredentials[*TokenCredentials], error) {

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

	return &TokenCredentials{
		Claims: claims,
	}, nil
}
