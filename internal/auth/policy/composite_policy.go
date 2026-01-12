package policy

import (
	"github.com/v2code/b16/internal/auth"
)

type CompositePolicy[T any] struct {
	policies []auth.Policy[T]
}

func NewCompositePolicy[T any](policies ...auth.Policy[T]) auth.Policy[T] {
	return &CompositePolicy[T]{policies: policies}
}

func (a *CompositePolicy[T]) Check(principal auth.Principal[T]) error {

	for _, p := range a.policies {
		if err := p.Check(principal); err != nil {
			return auth.ErrForbidden
		}
	}

	return nil
}
