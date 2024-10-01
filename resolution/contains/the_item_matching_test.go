package contains_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/resolution/contains"
	"github.com/grandper/go-screenplay/resolution/testdata"
)

func TestContainsTheItemMatchingResolution(t *testing.T) {
	t.Parallel()

	const regexStr = `H\w{4}`

	regex := regexp.MustCompile(regexStr)
	matcher1 := contains.TheItemMatching(regex)
	matcher2 := contains.TheItemMatchingRegexString(regexStr)

	createChannel := func(texts ...string) <-chan string {
		out := make(chan string)

		go func() {
			for _, text := range texts {
				out <- text
			}

			close(out)
		}()

		return out
	}

	t.Run("should match if we pass a matching value", func(t *testing.T) {
		t.Parallel()
		testdata.AssertMatch(t, matcher1, "Hello world!")
		testdata.AssertMatch(t, matcher2, "Hello world!")
	})

	t.Run("should match if the iterable contains the value", func(t *testing.T) {
		t.Parallel()
		testdata.AssertMatch(t, matcher1, createChannel("a", "Hello World!", "b"))
		testdata.AssertMatch(t, matcher1, []string{"Hello World!"})

		slice := []string{"Hello World!"}
		testdata.AssertMatch(t, matcher1, &slice)
		testdata.AssertMatch(t, matcher1, []string{"a", "Hello World", "b"})
	})

	t.Run("fails when the object is nil", func(t *testing.T) {
		t.Parallel()
		testdata.AssertNoMatch(t, matcher1, nil)

		var slice *[]string
		testdata.AssertNoMatch(t, matcher1, slice)
	})

	t.Run("fails when the iterable is empty", func(t *testing.T) {
		t.Parallel()
		testdata.AssertNoMatch(t, matcher1, []string{})
	})

	t.Run("fails when the iterable doesn't contains the item", func(t *testing.T) {
		t.Parallel()
		testdata.AssertNoMatch(t, matcher1, createChannel("a", "b", "c"))

		slice := []string{"a", "b"}
		testdata.AssertNoMatch(t, matcher1, &slice)
	})

	t.Run("returns an error when the map are used", func(t *testing.T) {
		t.Parallel()
		testdata.AssertMatcherFails(t, matcher1, map[string]int{"a": 1, "b": 2, "c": 3})
		testdata.AssertMatcherFails(t, matcher2, map[string]int{"a": 1, "b": 2, "c": 3})
	})

	t.Run("returns an error when the value is not a string", func(t *testing.T) {
		t.Parallel()
		testdata.AssertMatcherFails(t, matcher1, 1)
		testdata.AssertMatcherFails(t, matcher1, []int{1, 2})
	})

	t.Run("returns an error when passing a slice of pointers to a non-string type", func(t *testing.T) {
		t.Parallel()

		num := 42
		ptrSlice := []*int{&num}
		testdata.AssertMatcherFails(t, matcher1, ptrSlice)
	})

	t.Run("returns an error when passing channel of slice of pointers of strings with nil items", func(t *testing.T) {
		t.Parallel()

		ptrSlice := []*string{nil}

		ch := make(chan []*string, 1)
		ch <- ptrSlice

		testdata.AssertMatcherFails(t, matcher1, ch)
	})

	t.Run("returns an error when passing a slice of nil pointers", func(t *testing.T) {
		t.Parallel()

		ptrSlice := []*string{nil}
		testdata.AssertMatcherFails(t, matcher1, ptrSlice)
	})

	t.Run("returns an error when", func(t *testing.T) {
		t.Run("returns an error when the map are used", func(t *testing.T) {
			t.Parallel()
			testdata.AssertMatcherFails(t, matcher1, map[string]int{"a": 1, "b": 2, "c": 3})
			testdata.AssertMatcherFails(t, matcher2, map[string]int{"a": 1, "b": 2, "c": 3})
		})
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, `containing an item which matches pattern H\w{4}`, matcher1.String())
		assert.Equal(t, `containing an item which matches pattern H\w{4}`, matcher2.String())
	})
}
