package policy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/v2code/b16/internal/auth"
	"github.com/v2code/b16/internal/domain"
)

type TestPolicyParams struct {
	Name      string
	Policies  []domain.Policy[*dummyPrincipal]
	ExpectErr error
}

func TestAnyPolicy_Check(t *testing.T) {

	success := &fakePolicy[*dummyPrincipal]{err: nil}
	fail := &fakePolicy[*dummyPrincipal]{err: auth.ErrForbidden}

	tests := []TestPolicyParams{
		{
			Name:      "one policy succeeds",
			Policies:  []domain.Policy[*dummyPrincipal]{success},
			ExpectErr: nil,
		},
		{
			Name:      "first fails second succeeds",
			Policies:  []domain.Policy[*dummyPrincipal]{fail, success},
			ExpectErr: nil,
		},
		{
			Name:      "all policies fail",
			Policies:  []domain.Policy[*dummyPrincipal]{fail, fail},
			ExpectErr: auth.ErrForbidden,
		},
		{
			Name:      "no policies configured",
			Policies:  []domain.Policy[*dummyPrincipal]{},
			ExpectErr: auth.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			policy := NewAnyPolicy(tt.Policies...)

			err := policy.Check(&dummyPrincipal{})

			if tt.ExpectErr != nil {
				assert.ErrorIs(t, err, tt.ExpectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
