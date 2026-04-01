package screenplay_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/screenplay"
)

func TestResolutionToSeeThatTheObject(t *testing.T) {
	t.Run("can be created", func(t *testing.T) {
		t.Run("with a description and a successful matcher", func(t *testing.T) {
			t.Parallel()

			resolution := screenplay.ResolutionToSeeThatTheObject("value is 42", func() screenplay.Matcher {
				return func(obj any) (bool, error) {
					return obj == 42, nil
				}
			})

			assert.Equal(t, "value is 42", resolution.String())

			matcher := resolution.Resolve()
			ok, err := matcher(42)
			require.NoError(t, err)
			assert.True(t, ok)

			ok, err = matcher(0)
			require.NoError(t, err)
			assert.False(t, ok)
		})

		t.Run("with a description and a failing matcher", func(t *testing.T) {
			t.Parallel()

			expectedErr := errors.New("matcher failed")
			resolution := screenplay.ResolutionToSeeThatTheObject("is valid", func() screenplay.Matcher {
				return func(_ any) (bool, error) {
					return false, expectedErr
				}
			})

			assert.Equal(t, "is valid", resolution.String())

			matcher := resolution.Resolve()
			ok, err := matcher(nil)
			require.ErrorIs(t, err, expectedErr)
			assert.False(t, ok)
		})
	})
}
