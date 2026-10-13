package sized_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/internal/sized"
)

func TestSized(t *testing.T) {

	t.Run("len", func(t *testing.T) {
		t.Parallel()

		for _, testCase := range []struct {
			name    string
			value   any
			length  int
			isSized bool
		}{
			{"a nil value", nil, 0, true},
			{"a slice", []int{1, 2}, 2, true},
			{"an array", [3]int{}, 3, true},
			{"a map", map[string]int{"a": 1}, 1, true},
			{"a string", "hello", 5, true},
			{"a pointer to a slice", &[]int{1}, 1, true},
			{"a nil pointer to a slice", (*[]int)(nil), 0, true},
			{"a number", 3, 0, false},
			{"a pointer to a number", new(int), 0, false},
		} {
			t.Run(testCase.name, func(t *testing.T) {
				t.Parallel()

				length, isSized := sized.Len(testCase.value)
				assert.Equal(t, testCase.length, length)
				assert.Equal(t, testCase.isSized, isSized)
			})
		}
	})

	t.Run("check", func(t *testing.T) {
		t.Parallel()

		t.Run("accepts the types whose values have a length", func(t *testing.T) {
			t.Parallel()
			require.NoError(t, sized.Check[[]int]())
			require.NoError(t, sized.Check[string]())
			require.NoError(t, sized.Check[map[string]int]())
			require.NoError(t, sized.Check[chan int]())
			require.NoError(t, sized.Check[*[]int]())
			require.NoError(t, sized.Check[any]())
		})

		t.Run("rejects the types whose values never have a length", func(t *testing.T) {
			t.Parallel()
			require.EqualError(t, sized.Check[int](), "values of type int have no length")
			require.Error(t, sized.Check[*int]())
		})
	})
}
