package policy

import (
	"slices"

	"github.com/v2code/b16/internal/auth"
	"github.com/v2code/b16/internal/auth/manager"
	"github.com/v2code/b16/internal/domain"
)

type RequireRole struct {
	roles []string
}

func RequireRolePolicy(roles ...string) domain.Policy[*manager.TokenPrincipal] {
	return &RequireRole{roles: roles}
}

func (p *RequireRole) Check(credentials domain.Principal[*manager.TokenPrincipal]) error {

	for _, role := range p.roles {
		if !slices.Contains(credentials.Principal().Roles, role) {
			return auth.ErrForbidden
		}
	}

	return nil
}
