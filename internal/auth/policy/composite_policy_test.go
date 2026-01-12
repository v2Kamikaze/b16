package policy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/v2code/b16/internal/auth"
)

func TestCompositePolicy_Check(t *testing.T) {

	success := &fakePolicy[*dummyPrincipal]{err: nil}
	fail := &fakePolicy[*dummyPrincipal]{err: auth.ErrForbidden}

	tests := []TestPolicyParams{
		{
			Name:      "all policy succeeds",
			Policies:  []auth.Policy[*dummyPrincipal]{success},
			ExpectErr: nil,
		},
		{
			Name:      "one policy fails",
			Policies:  []auth.Policy[*dummyPrincipal]{success, fail},
			ExpectErr: auth.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			policy := NewCompositePolicy(tt.Policies...)

			err := policy.Check(&dummyPrincipal{})

			if tt.ExpectErr != nil {
				assert.ErrorIs(t, err, tt.ExpectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
