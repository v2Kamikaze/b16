package manager

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/v2code/b16/internal/auth"
)

func TestBasicAuthManager_Authenticate(t *testing.T) {
	manager := NewBasicAuthManager(BasicAuthParams{
		Username: "admin",
		Password: "secret",
	})

	tests := []struct {
		name         string
		setupRequest func(r *http.Request)
		expectError  bool
		expectedUser string
		expectedPass string
	}{
		{
			name: "no basic auth header",
			setupRequest: func(r *http.Request) {
				// nada
			},
			expectError: true,
		},
		{
			name: "invalid credentials",
			setupRequest: func(r *http.Request) {
				r.SetBasicAuth("admin", "wrong")
			},
			expectError: true,
		},
		{
			name: "valid credentials",
			setupRequest: func(r *http.Request) {
				r.SetBasicAuth("admin", "secret")
			},
			expectError:  false,
			expectedUser: "admin",
			expectedPass: "secret",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			tt.setupRequest(req)

			cred, err := manager.Authenticate(req)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if err != auth.ErrUnauthorized {
					t.Fatalf("expected ErrUnauthorized, got %v", err)
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

			if principal.Username != tt.expectedUser {
				t.Errorf("expected username %q, got %q", tt.expectedUser, principal.Username)
			}

			if principal.Password != tt.expectedPass {
				t.Errorf("expected password %q, got %q", tt.expectedPass, principal.Password)
			}
		})
	}
}
