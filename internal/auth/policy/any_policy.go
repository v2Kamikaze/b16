package policy

import (
	"github.com/v2code/b16/internal/auth"
)

type AnyPolicy[T any] struct {
	policies []auth.Policy[T]
}

func NewAnyPolicy[T any](policies ...auth.Policy[T]) auth.Policy[T] {
	return &AnyPolicy[T]{policies: policies}
}

func (a *AnyPolicy[T]) Check(principal auth.Principal[T]) error {
	for _, p := range a.policies {
		if err := p.Check(principal); err == nil {
			return nil
		}
	}

	return auth.ErrForbidden
}
