package length

import (
	"errors"
	"fmt"
	"math"

	"github.com/grandper/go-screenplay/v2/internal/verdict"
	"github.com/grandper/go-screenplay/v2/resolution"
	"github.com/grandper/go-screenplay/v2/resolution/internal/sized"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// Exactly creates a resolution to tell if a value has exactly the given length.
func Exactly[T any](length int) screenplay.Resolution[T] {
	return within[T](fmt.Sprintf("length exactly %d", length), length, length)
}

// AtLeast creates a resolution to tell if a value has at least the given length.
func AtLeast[T any](shortest int) screenplay.Resolution[T] {
	return within[T](fmt.Sprintf("length at least %d", shortest), shortest, math.MaxInt)
}

// AtMost creates a resolution to tell if a value has at most the given length.
func AtMost[T any](longest int) screenplay.Resolution[T] {
	return within[T](fmt.Sprintf("length at most %d", longest), 0, longest)
}

// Between creates a resolution to tell if the length of a value lies between
// two bounds, both included.
func Between[T any](shortest, longest int) screenplay.Resolution[T] {
	return within[T](fmt.Sprintf("length between %d and %d", shortest, longest), shortest, longest)
}

// within creates a resolution to tell if the length of a value lies between
// two bounds. The resolution is invalid when values of type T can never have a
// length, when a bound is negative, or when the bounds are out of order.
func within[T any](description string, shortest, longest int) screenplay.Resolution[T] {
	if err := sized.Check[T](); err != nil {
		return verdict.Invalid[T](description, err)
	}

	if shortest < 0 || longest < 0 {
		return verdict.Invalid[T](description, errors.New("a length must not be negative"))
	}

	if shortest > longest {
		return verdict.Invalid[T](description, errors.New("the shortest length is greater than the longest length"))
	}

	return resolution.FromFunc(
		description,
		func(value T) error {
			length, isSized := sized.Len(value)

			return verdict.Of(isSized && length >= shortest && length <= longest, value)
		},
	)
}
