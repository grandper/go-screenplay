package contain_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/resolution/contain"
	"github.com/grandper/go-screenplay/v2/resolution/fixture"
)

func TestContain(t *testing.T) {
	t.Parallel()

	t.Run("substring", func(t *testing.T) {
		t.Parallel()

		resolution := contain.Substring("world")

		fixture.AssertMatch(t, resolution, "hello world")
		fixture.AssertNoMatch(t, resolution, "hello everybody")
		require.EqualError(t, resolution.Inspect("hi"), `got "hi"`)
		assert.Equal(t, `contain substring "world"`, resolution.String())
	})

	t.Run("the bytes", func(t *testing.T) {
		t.Parallel()

		resolution := contain.TheBytes([]byte("lo Wor"))

		fixture.AssertMatch(t, resolution, []byte("Hello World"))
		fixture.AssertNoMatch(t, resolution, []byte("Hello"))
		fixture.AssertNoMatch(t, resolution, nil)
		assert.Equal(t, `contain the bytes "lo Wor"`, resolution.String())

		t.Run("is not affected by a later change of the given bytes", func(t *testing.T) {
			t.Parallel()

			given := []byte("abc")
			containABC := contain.TheBytes(given)
			given[0] = 'x'

			fixture.AssertMatch(t, containABC, []byte("abc"))
		})
	})

	t.Run("element", func(t *testing.T) {
		t.Parallel()

		resolution := contain.Element("Adam")

		fixture.AssertMatch(t, resolution, []string{"Adam", "Eve"})
		fixture.AssertNoMatch(t, resolution, []string{"Eve"})
		fixture.AssertNoMatch(t, resolution, nil)
		fixture.AssertMatch(t, contain.Element(3), []int{1, 2, 3})
		assert.Equal(t, `contain element "Adam"`, resolution.String())
	})

	t.Run("key", func(t *testing.T) {
		t.Parallel()

		resolution := contain.Key[string, int]("total")

		fixture.AssertMatch(t, resolution, map[string]int{"total": 0})
		fixture.AssertNoMatch(t, resolution, map[string]int{"count": 1})
		fixture.AssertNoMatch(t, resolution, nil)
		assert.Equal(t, `contain key "total"`, resolution.String())
	})

	t.Run("value", func(t *testing.T) {
		t.Parallel()

		resolution := contain.Value[string, int](42)

		fixture.AssertMatch(t, resolution, map[string]int{"total": 42})
		fixture.AssertNoMatch(t, resolution, map[string]int{"total": 1})
		fixture.AssertNoMatch(t, resolution, nil)
		assert.Equal(t, "contain value 42", resolution.String())
	})

	t.Run("entries", func(t *testing.T) {
		t.Parallel()

		resolution := contain.Entries(map[string]string{"Content-Type": "application/json", "Accept": "*/*"})

		t.Run("is satisfied by a map holding every entry", func(t *testing.T) {
			t.Parallel()
			fixture.AssertMatch(t, resolution, map[string]string{
				"Content-Type": "application/json",
				"Accept":       "*/*",
				"Host":         "example.com",
			})
		})

		t.Run("is not satisfied when an entry is missing or differs", func(t *testing.T) {
			t.Parallel()
			fixture.AssertNoMatch(t, resolution, map[string]string{"Content-Type": "application/json"})
			fixture.AssertNoMatch(t, resolution, map[string]string{"Content-Type": "text/html", "Accept": "*/*"})
			fixture.AssertNoMatch(t, resolution, nil)
		})

		t.Run("describes itself with its entries in a stable order", func(t *testing.T) {
			t.Parallel()
			assert.Equal(t,
				`contain entries "Accept": "*/*", "Content-Type": "application/json"`,
				resolution.String(),
			)
		})

		t.Run("is not affected by a later change of the given map", func(t *testing.T) {
			t.Parallel()

			given := map[string]int{"a": 1}
			containA := contain.Entries(given)
			given["b"] = 2

			fixture.AssertMatch(t, containA, map[string]int{"a": 1})
		})
	})

	t.Run("the item matching", func(t *testing.T) {
		t.Parallel()

		resolution := contain.AnItemMatching(regexp.MustCompile(`^A-\d+$`))

		fixture.AssertMatch(t, resolution, []string{"B-1", "A-42"})
		fixture.AssertNoMatch(t, resolution, []string{"B-1", "A-"})
		fixture.AssertNoMatch(t, resolution, nil)
		assert.Equal(t, "contain an item matching pattern `^A-\\d+$`", resolution.String())

		t.Run("accepts the pattern as a text", func(t *testing.T) {
			t.Parallel()
			fixture.AssertMatch(t, contain.AnItemMatchingPatternString(`^A-\d+$`), []string{"A-42"})
		})

		t.Run("is invalid when the pattern is nil or does not compile", func(t *testing.T) {
			t.Parallel()
			fixture.AssertInvalid(t, contain.AnItemMatching(nil), []string{"A-42"})
			fixture.AssertInvalid(t, contain.AnItemMatchingPatternString("A-("), []string{"A-42"})
		})
	})

	t.Run("with an interface type parameter", func(t *testing.T) {
		fixture.AssertMatch(t, contain.Element[any](3), []any{1, 3})
		fixture.AssertInvalid(t, contain.Element[any]([]int{1}), []any{[]int{1}})

		fixture.AssertMatch(t, contain.Value[string, any](3), map[string]any{"a": 3})
		fixture.AssertInvalid(t, contain.Value[string, any]([]int{1}), map[string]any{"a": []int{1}})

		fixture.AssertMatch(t, contain.Entries[string, any](map[string]any{"a": 3}), map[string]any{"a": 3, "b": 4})
		fixture.AssertInvalid(
			t,
			contain.Entries[string, any](map[string]any{"a": []int{1}}),
			map[string]any{"a": []int{1}},
		)
	})
}
