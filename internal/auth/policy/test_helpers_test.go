package policy

import (
	"github.com/v2code/b16/internal/domain"
)

type dummyPrincipal struct{}

func (d *dummyPrincipal) Principal() *dummyPrincipal {
	return d
}

type fakePolicy[T any] struct {
	err error
}

func (p *fakePolicy[T]) Check(credentials domain.Principal[T]) error {
	return p.err
}
