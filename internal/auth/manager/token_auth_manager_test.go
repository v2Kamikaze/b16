package manager

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/v2code/b16/internal/auth"
	"github.com/v2code/b16/internal/security"
)

type fakeTokenIssuer struct {
	claims *security.Claims
	err    error
}

func (f *fakeTokenIssuer) Create(data *security.Claims) (string, error) {
	return "fake-token", nil
}

func (f *fakeTokenIssuer) Decode(token string) (*security.Claims, error) {
	return f.claims, f.err
}

type TokenAuthTestParams struct {
	Name        string
	AuthHeader  string
	Issuer      *fakeTokenIssuer
	ExpectErr   error
	ExpectEmail string
	ExpectRole  string
}

func TestTokenAuthManager_Authenticate(t *testing.T) {
	validClaims := &security.Claims{
		Email: "admin@email.com",
		Roles: []string{"ADMIN"},
	}

	cases := []TokenAuthTestParams{
		{
			Name:       "missing authorization header",
			AuthHeader: "",
			Issuer:     &fakeTokenIssuer{},
			ExpectErr:  auth.ErrUnauthorized,
		},
		{
			Name:       "empty bearer token",
			AuthHeader: "Bearer ",
			Issuer:     &fakeTokenIssuer{},
			ExpectErr:  auth.ErrTokenMissing,
		},
		{
			Name:       "invalid token",
			AuthHeader: "Bearer invalid-token",
			Issuer: &fakeTokenIssuer{
				err: auth.ErrUnauthorized,
			},
			ExpectErr: auth.ErrUnauthorized,
		},
		{
			Name:       "valid token",
			AuthHeader: "Bearer valid-token",
			Issuer: &fakeTokenIssuer{
				claims: validClaims,
			},
			ExpectEmail: "admin@email.com",
			ExpectRole:  "ADMIN",
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			manager := NewTokenAuthManager(tt.Issuer)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.AuthHeader != "" {
				req.Header.Set("Authorization", tt.AuthHeader)
			}

			cred, err := manager.Authenticate(req)

			if tt.ExpectErr != nil {
				require.ErrorIs(t, err, tt.ExpectErr)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, cred)

			principal := cred.Principal()
			require.Equal(t, tt.ExpectEmail, principal.Email)
			require.Contains(t, principal.Roles, tt.ExpectRole)
		})
	}
}
