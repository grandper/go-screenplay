package contain_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/resolution/contain"
	"github.com/grandper/go-screenplay/v2/resolution/testdata"
)

func TestContain(t *testing.T) {

	t.Run("substring", func(t *testing.T) {
		t.Parallel()

		resolution := contain.Substring("world")

		testdata.AssertMatch(t, resolution, "hello world")
		testdata.AssertNoMatch(t, resolution, "hello everybody")
		require.EqualError(t, resolution.Inspect("hi"), `got "hi"`)
		assert.Equal(t, `contain substring "world"`, resolution.String())
	})

	t.Run("the bytes", func(t *testing.T) {
		t.Parallel()

		resolution := contain.TheBytes([]byte("lo Wor"))

		testdata.AssertMatch(t, resolution, []byte("Hello World"))
		testdata.AssertNoMatch(t, resolution, []byte("Hello"))
		testdata.AssertNoMatch(t, resolution, nil)
		assert.Equal(t, `contain the bytes "lo Wor"`, resolution.String())

		t.Run("is not affected by a later change of the given bytes", func(t *testing.T) {
			t.Parallel()

			given := []byte("abc")
			containABC := contain.TheBytes(given)
			given[0] = 'x'

			testdata.AssertMatch(t, containABC, []byte("abc"))
		})
	})

	t.Run("element", func(t *testing.T) {
		t.Parallel()

		resolution := contain.Element("Adam")

		testdata.AssertMatch(t, resolution, []string{"Adam", "Eve"})
		testdata.AssertNoMatch(t, resolution, []string{"Eve"})
		testdata.AssertNoMatch(t, resolution, nil)
		testdata.AssertMatch(t, contain.Element(3), []int{1, 2, 3})
		assert.Equal(t, `contain element "Adam"`, resolution.String())
	})

	t.Run("key", func(t *testing.T) {
		t.Parallel()

		resolution := contain.Key[string, int]("total")

		testdata.AssertMatch(t, resolution, map[string]int{"total": 0})
		testdata.AssertNoMatch(t, resolution, map[string]int{"count": 1})
		testdata.AssertNoMatch(t, resolution, nil)
		assert.Equal(t, `contain key "total"`, resolution.String())
	})

	t.Run("value", func(t *testing.T) {
		t.Parallel()

		resolution := contain.Value[string, int](42)

		testdata.AssertMatch(t, resolution, map[string]int{"total": 42})
		testdata.AssertNoMatch(t, resolution, map[string]int{"total": 1})
		testdata.AssertNoMatch(t, resolution, nil)
		assert.Equal(t, "contain value 42", resolution.String())
	})

	t.Run("entries", func(t *testing.T) {
		t.Parallel()

		resolution := contain.Entries(map[string]string{"Content-Type": "application/json", "Accept": "*/*"})

		t.Run("is satisfied by a map holding every entry", func(t *testing.T) {
			t.Parallel()
			testdata.AssertMatch(t, resolution, map[string]string{
				"Content-Type": "application/json",
				"Accept":       "*/*",
				"Host":         "example.com",
			})
		})

		t.Run("is not satisfied when an entry is missing or differs", func(t *testing.T) {
			t.Parallel()
			testdata.AssertNoMatch(t, resolution, map[string]string{"Content-Type": "application/json"})
			testdata.AssertNoMatch(t, resolution, map[string]string{"Content-Type": "text/html", "Accept": "*/*"})
			testdata.AssertNoMatch(t, resolution, nil)
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

			testdata.AssertMatch(t, containA, map[string]int{"a": 1})
		})
	})

	t.Run("the item matching", func(t *testing.T) {
		t.Parallel()

		resolution := contain.TheItemMatching(regexp.MustCompile(`^A-\d+$`))

		testdata.AssertMatch(t, resolution, []string{"B-1", "A-42"})
		testdata.AssertNoMatch(t, resolution, []string{"B-1", "A-"})
		testdata.AssertNoMatch(t, resolution, nil)
		assert.Equal(t, "contain an item matching pattern `^A-\\d+$`", resolution.String())

		t.Run("accepts the pattern as a text", func(t *testing.T) {
			t.Parallel()
			testdata.AssertMatch(t, contain.TheItemMatchingRegexString(`^A-\d+$`), []string{"A-42"})
		})

		t.Run("is invalid when the pattern is nil or does not compile", func(t *testing.T) {
			t.Parallel()
			testdata.AssertInvalid(t, contain.TheItemMatching(nil), []string{"A-42"})
			testdata.AssertInvalid(t, contain.TheItemMatchingRegexString("A-("), []string{"A-42"})
		})
	})

	t.Run("with an interface type parameter", func(t *testing.T) {
		testdata.AssertMatch(t, contain.Element[any](3), []any{1, 3})
		testdata.AssertInvalid(t, contain.Element[any]([]int{1}), []any{[]int{1}})

		testdata.AssertMatch(t, contain.Value[string, any](3), map[string]any{"a": 3})
		testdata.AssertInvalid(t, contain.Value[string, any]([]int{1}), map[string]any{"a": []int{1}})

		testdata.AssertMatch(t, contain.Entries[string, any](map[string]any{"a": 3}), map[string]any{"a": 3, "b": 4})
		testdata.AssertInvalid(t, contain.Entries[string, any](map[string]any{"a": []int{1}}), map[string]any{"a": []int{1}})
	})
}
