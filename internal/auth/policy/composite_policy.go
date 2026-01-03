package policy

import "github.com/v2code/b16/internal/domain"

type CompositePolicy[T any] struct {
	policies []domain.Policy[T]
}

func NewCompositePolicy[T any](policies ...domain.Policy[T]) domain.Policy[T] {

	if len(policies) == 0 {
		panic("composite policy requires at least one policy")
	}

	return &CompositePolicy[T]{policies: policies}
}

func (a *CompositePolicy[T]) Check(cred domain.UserCredentials[T]) error {
	for _, p := range a.policies {
		if err := p.Check(cred); err != nil {
			return err
		}
	}
	return nil
}
