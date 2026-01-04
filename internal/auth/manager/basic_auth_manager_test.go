package manager

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/v2code/b16/internal/auth"
)

type TestBasicAuthParams struct {
	Name         string
	SetupRequest func(r *http.Request)
	ExpectedErr  error
	ExpectedUser string
	ExpectedPass string
}

func TestBasicAuthManager_Authenticate(t *testing.T) {
	manager := NewBasicAuthManager(map[string]string{
		"admin": "secret",
		"user":  "password",
	})

	cases := []TestBasicAuthParams{
		{
			Name:         "no basic auth header",
			SetupRequest: func(r *http.Request) {},
			ExpectedErr:  auth.ErrUnauthorized,
		},
		{
			Name: "invalid credentials",
			SetupRequest: func(r *http.Request) {
				r.SetBasicAuth("admin", "wrong")
			},
			ExpectedErr: auth.ErrUnauthorized,
		},
		{
			Name: "valid credentials - admin",
			SetupRequest: func(r *http.Request) {
				r.SetBasicAuth("admin", "secret")
			},
			ExpectedErr:  nil,
			ExpectedUser: "admin",
			ExpectedPass: "secret",
		},
		{
			Name: "valid credentials - user",
			SetupRequest: func(r *http.Request) {
				r.SetBasicAuth("user", "password")
			},
			ExpectedErr:  nil,
			ExpectedUser: "user",
			ExpectedPass: "password",
		},
		{
			Name: "non-existent user",
			SetupRequest: func(r *http.Request) {
				r.SetBasicAuth("unknown", "password")
			},
			ExpectedErr: auth.ErrUnauthorized,
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			tt.SetupRequest(req)

			cred, err := manager.Authenticate(req)

			if tt.ExpectedErr != nil {
				require.ErrorIs(t, err, tt.ExpectedErr)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, cred)

			principal := cred.Principal()
			require.Equal(t, tt.ExpectedUser, principal.Username)
			require.Equal(t, tt.ExpectedPass, principal.Password)
		})
	}
}
