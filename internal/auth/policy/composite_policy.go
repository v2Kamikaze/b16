package policy

import (
	"github.com/v2code/b16/internal/auth"
	"github.com/v2code/b16/internal/domain"
)

type CompositePolicy[T any] struct {
	policies []domain.Policy[T]
}

func NewCompositePolicy[T any](policies ...domain.Policy[T]) domain.Policy[T] {
	return &CompositePolicy[T]{policies: policies}
}

func (a *CompositePolicy[T]) Check(credentials domain.Principal[T]) error {

	for _, p := range a.policies {
		if err := p.Check(credentials); err != nil {
			return auth.ErrForbidden
		}
	}

	return nil
}
