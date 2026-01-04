package policy

import (
	"github.com/v2code/b16/internal/auth"
	"github.com/v2code/b16/internal/domain"
)

type AnyPolicy[T any] struct {
	policies []domain.Policy[T]
}

func NewAnyPolicy[T any](policies ...domain.Policy[T]) domain.Policy[T] {
	return &AnyPolicy[T]{policies: policies}
}

func (a *AnyPolicy[T]) Check(credentials domain.Principal[T]) error {
	for _, p := range a.policies {
		if err := p.Check(credentials); err == nil {
			return nil
		}
	}

	return auth.ErrForbidden
}
