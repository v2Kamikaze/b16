package policy

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/v2code/b16/internal/auth"
	"github.com/v2code/b16/internal/auth/manager"
	"github.com/v2code/b16/internal/security"
)

type TestRequireRolePolicyParams struct {
	Name        string
	PolicyRoles []string
	UserRoles   []string
	ExpectErr   error
}

func TestRequireRolePolicy_Check(t *testing.T) {
	cases := []TestRequireRolePolicyParams{
		{Name: "user has all required roles",
			PolicyRoles: []string{"ADMIN", "USER"},
			UserRoles:   []string{"ADMIN", "USER", "AUDITOR"},
			ExpectErr:   nil,
		},
		{Name: "user missing one required role",
			PolicyRoles: []string{"ADMIN", "USER"},
			UserRoles:   []string{"ADMIN"},
			ExpectErr:   auth.ErrForbidden,
		},
		{Name: "user missing all required roles",
			PolicyRoles: []string{"ADMIN"},
			UserRoles:   []string{"USER"},
			ExpectErr:   auth.ErrForbidden,
		},
		{Name: "single required role and user has it",
			PolicyRoles: []string{"ADMIN"},
			UserRoles:   []string{"ADMIN"},
			ExpectErr:   nil,
		},
		{Name: "user has no roles",
			PolicyRoles: []string{"ADMIN", "USER"},
			UserRoles:   []string{},
			ExpectErr:   auth.ErrForbidden,
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			policy := RequireRolePolicy(tt.PolicyRoles...)

			credentials := &manager.TokenPrincipal{
				Claims: &security.Claims{
					Roles: tt.UserRoles,
				},
			}

			err := policy.Check(credentials)

			if tt.ExpectErr != nil {
				require.ErrorIs(t, err, tt.ExpectErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
