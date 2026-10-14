package contain

import (
	"bytes"
	"errors"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/grandper/go-screenplay/v2/internal/verdict"
	"github.com/grandper/go-screenplay/v2/resolution"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// Substring creates a resolution to tell if a text contains a given text.
func Substring(substring string) screenplay.Resolution[string] {
	return resolution.FromFunc(
		"contain substring "+verdict.Describe(substring),
		func(text string) error {
			return verdict.Of(strings.Contains(text, substring), text)
		},
	)
}

// TheBytes creates a resolution to tell if a slice of bytes contains a given slice of bytes.
func TheBytes(subslice []byte) screenplay.Resolution[[]byte] {
	subslice = bytes.Clone(subslice)

	return resolution.FromFunc(
		"contain the bytes "+verdict.Describe(subslice),
		func(value []byte) error {
			return verdict.Of(bytes.Contains(value, subslice), value)
		},
	)
}

// Element creates a resolution to tell if a slice contains a given element.
// The type of the slice is inferred from the element.
func Element[T comparable](element T) screenplay.Resolution[[]T] {
	return resolution.FromFunc(
		"contain element "+verdict.Describe(element),
		func(elements []T) error {
			for _, candidate := range elements {
				equal, err := verdict.Equal(candidate, element)
				if err != nil {
					return err
				}
				if equal {
					return nil
				}
			}

			return verdict.Of(false, elements)
		},
	)
}

// Key creates a resolution to tell if a map contains a given key. Go cannot
// infer the type of the values of the map from the key alone, so both types are
// given explicitly:
//
//	contain.Key[string, int]("total")
func Key[K comparable, V any](key K) screenplay.Resolution[map[K]V] {
	return resolution.FromFunc(
		"contain key "+verdict.Describe(key),
		func(entries map[K]V) error {
			_, found := entries[key]

			return verdict.Of(found, entries)
		},
	)
}

// Value creates a resolution to tell if a map contains a given value. Go cannot
// infer the type of the keys of the map from the value alone, so both types are
// given explicitly:
//
//	contain.Value[string, int](42)
func Value[K comparable, V comparable](value V) screenplay.Resolution[map[K]V] {
	return resolution.FromFunc(
		"contain value "+verdict.Describe(value),
		func(entries map[K]V) error {
			for _, candidate := range entries {
				equal, err := verdict.Equal(candidate, value)
				if err != nil {
					return err
				}
				if equal {
					return nil
				}
			}

			return verdict.Of(false, entries)
		},
	)
}

// Entries creates a resolution to tell if a map contains every entry of a given
// map: each of its keys, associated with the same value.
func Entries[K comparable, V comparable](expected map[K]V) screenplay.Resolution[map[K]V] {
	expected = maps.Clone(expected)

	return resolution.FromFunc(
		"contain entries "+describeEntries(expected),
		func(entries map[K]V) error {
			for key, value := range expected {
				candidate, found := entries[key]
				if !found {
					return verdict.Of(false, entries)
				}
				equal, err := verdict.Equal(candidate, value)
				if err != nil {
					return err
				}
				if !equal {
					return verdict.Of(false, entries)
				}
			}

			return nil
		},
	)
}

// describeEntries renders the entries of a map as "key: value" pairs.
func describeEntries[K comparable, V comparable](entries map[K]V) string {
	described := make([]string, 0, len(entries))
	for key, value := range entries {
		described = append(described, verdict.Describe(key)+": "+verdict.Describe(value))
	}

	slices.Sort(described)

	return strings.Join(described, ", ")
}

// AnItemMatching creates a resolution to tell if a slice of texts contains an
// item matching a regular expression. The resolution is invalid when the
// expression is nil.
func AnItemMatching(pattern *regexp.Regexp) screenplay.Resolution[[]string] {
	if pattern == nil {
		return verdict.Invalid[[]string](
			"contain an item matching pattern <nil>",
			errors.New("the regular expression is nil"),
		)
	}

	return resolution.FromFunc(
		fmt.Sprintf("contain an item matching pattern `%s`", pattern),
		func(items []string) error {
			return verdict.Of(slices.ContainsFunc(items, pattern.MatchString), items)
		},
	)
}

// AnItemMatchingPatternString creates a resolution to tell if a slice of texts
// contains an item matching a regular expression given as a text. The
// resolution is invalid when the expression does not compile.
func AnItemMatchingPatternString(pattern string) screenplay.Resolution[[]string] {
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return verdict.Invalid[[]string](fmt.Sprintf("contain an item matching pattern `%s`", pattern), err)
	}

	return AnItemMatching(compiled)
}
