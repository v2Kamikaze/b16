package manager

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestTokenAuthManager_Authenticate(t *testing.T) {

	validClaims := &security.Claims{
		Email: "admin@email.com",
		Roles: []string{"ADMIN"},
	}

	tests := []struct {
		name        string
		authHeader  string
		issuer      *fakeTokenIssuer
		expectError error
		expectEmail string
		expectRole  string
	}{
		{
			name:        "missing authorization header",
			authHeader:  "",
			issuer:      &fakeTokenIssuer{},
			expectError: auth.ErrUnauthorized,
		},
		{
			name:        "empty bearer token",
			authHeader:  "Bearer ",
			issuer:      &fakeTokenIssuer{},
			expectError: auth.ErrTokenMissing,
		},
		{
			name:       "invalid token",
			authHeader: "Bearer invalid-token",
			issuer: &fakeTokenIssuer{
				err: auth.ErrUnauthorized,
			},
			expectError: auth.ErrUnauthorized,
		},
		{
			name:       "valid token",
			authHeader: "Bearer valid-token",
			issuer: &fakeTokenIssuer{
				claims: validClaims,
			},
			expectEmail: "admin@email.com",
			expectRole:  "ADMIN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			manager := NewTokenAuthManager(tt.issuer)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			cred, err := manager.Authenticate(req)

			if tt.expectError != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.expectError)
				}
				if err != tt.expectError {
					t.Fatalf("expected error %v, got %v", tt.expectError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cred == nil {
				t.Fatalf("expected credentials, got nil")
			}

			principal := cred.Principal()

			if principal.Email != tt.expectEmail {
				t.Errorf("expected email %q, got %q", tt.expectEmail, principal.Email)
			}

			if len(principal.Roles) == 0 || principal.Roles[0] != tt.expectRole {
				t.Errorf("expected role %q, got %v", tt.expectRole, principal.Roles)
			}
		})
	}
}
