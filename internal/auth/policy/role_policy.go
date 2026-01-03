package policy

import (
	"slices"

	"github.com/v2code/b16/internal/auth"
	"github.com/v2code/b16/internal/auth/manager"
	"github.com/v2code/b16/internal/domain"
)

type RequireRole struct {
	role string
}

func RequireRolePolicy(role string) domain.Policy[*manager.TokenCredentials] {
	return &RequireRole{role: role}
}

func (p *RequireRole) Check(cred domain.UserCredentials[*manager.TokenCredentials]) error {

	if slices.Contains(cred.GetCredentials().Roles, p.role) {
		return nil
	}

	return auth.ErrForbidden
}
