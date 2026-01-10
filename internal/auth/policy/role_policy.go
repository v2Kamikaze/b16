package policy

import (
	"slices"

	"github.com/v2code/b16/internal/auth"
	"github.com/v2code/b16/internal/auth/manager"
)

type RequireRole struct {
	roles []string
}

func RequireRolePolicy(roles ...string) auth.Policy[*manager.TokenPrincipal] {
	return &RequireRole{roles: roles}
}

func (p *RequireRole) Check(principal auth.Principal[*manager.TokenPrincipal]) error {

	for _, role := range p.roles {
		if !slices.Contains(principal.Principal().Roles, role) {
			return auth.ErrForbidden
		}
	}

	return nil
}
