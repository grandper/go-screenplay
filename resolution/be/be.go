package be

import (
	"cmp"
	"errors"
	"fmt"
	"math"
	"reflect"

	"github.com/grandper/go-screenplay/v2/internal/sized"
	"github.com/grandper/go-screenplay/v2/internal/verdict"
	"github.com/grandper/go-screenplay/v2/resolution"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// integer is the set of the integer types.
type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// float is the set of the floating-point types.
type float interface {
	~float32 | ~float64
}

// EqualTo creates a resolution to tell if a value equals a given value. The
// comparison is exact on the type: a question answering with an int64 needs
// be.EqualTo(int64(200)). Use DeepEqualTo for the values that are not
// comparable, such as slices and maps.
func EqualTo[T comparable](expected T) screenplay.Resolution[T] {
	return resolution.FromFunc(
		"be equal to "+verdict.Describe(expected),
		func(value T) error {
			return verdict.Of(value == expected, value)
		},
	)
}

// DeepEqualTo creates a resolution to tell if a value is deeply equal to a
// given value, as reflect.DeepEqual defines it. It is the resolution to
// compare slices, maps and structs holding them.
func DeepEqualTo[T any](expected T) screenplay.Resolution[T] {
	return resolution.FromFunc(
		"be deeply equal to "+verdict.Describe(expected),
		func(value T) error {
			return verdict.Of(reflect.DeepEqual(value, expected), value)
		},
	)
}

// GreaterThan creates a resolution to tell if a value is strictly greater than a bound.
func GreaterThan[T cmp.Ordered](bound T) screenplay.Resolution[T] {
	return resolution.FromFunc(
		"be greater than "+verdict.Describe(bound),
		func(value T) error {
			return verdict.Of(value > bound, value)
		},
	)
}

// AtLeast creates a resolution to tell if a value is greater than or equal to a bound.
func AtLeast[T cmp.Ordered](bound T) screenplay.Resolution[T] {
	return resolution.FromFunc(
		"be at least "+verdict.Describe(bound),
		func(value T) error {
			return verdict.Of(value >= bound, value)
		},
	)
}

// LessThan creates a resolution to tell if a value is strictly less than a bound.
func LessThan[T cmp.Ordered](bound T) screenplay.Resolution[T] {
	return resolution.FromFunc(
		"be less than "+verdict.Describe(bound),
		func(value T) error {
			return verdict.Of(value < bound, value)
		},
	)
}

// AtMost creates a resolution to tell if a value is less than or equal to a bound.
func AtMost[T cmp.Ordered](bound T) screenplay.Resolution[T] {
	return resolution.FromFunc(
		"be at most "+verdict.Describe(bound),
		func(value T) error {
			return verdict.Of(value <= bound, value)
		},
	)
}

// Between creates a resolution to tell if a value lies between two bounds,
// both included. The resolution is invalid when the lower bound is greater
// than the upper one.
func Between[T cmp.Ordered](lowerBound, upperBound T) screenplay.Resolution[T] {
	description := fmt.Sprintf("be between %s and %s", verdict.Describe(lowerBound), verdict.Describe(upperBound))

	if lowerBound > upperBound {
		return verdict.Invalid[T](description, errors.New("the lower bound is greater than the upper bound"))
	}

	return resolution.FromFunc(
		description,
		func(value T) error {
			return verdict.Of(value >= lowerBound && value <= upperBound, value)
		},
	)
}

// CloseTo creates a resolution to tell if a number is within a tolerance of a
// target, the bounds included. The resolution is invalid when the tolerance is
// negative or is not a number.
func CloseTo[T integer | float](target T, tolerance float64) screenplay.Resolution[T] {
	description := fmt.Sprintf("be close to %v within %v", target, tolerance)

	if tolerance < 0 || math.IsNaN(tolerance) {
		return verdict.Invalid[T](description, errors.New("the tolerance must be a positive number"))
	}

	return resolution.FromFunc(
		description,
		func(value T) error {
			return verdict.Of(math.Abs(float64(value)-float64(target)) <= tolerance, value)
		},
	)
}

// Empty creates a resolution to tell if a value of type T is empty: a channel,
// a map, a slice, an array or a string (or a pointer to one) of length zero.
// Go cannot infer T from the resolution alone, so it is given explicitly:
//
//	be.Empty[[]string]()
//
// The resolution is invalid when values of type T can never have a length.
func Empty[T any]() screenplay.Resolution[T] {
	if err := sized.Check[T](); err != nil {
		return verdict.Invalid[T]("be empty", err)
	}

	return resolution.FromFunc(
		"be empty",
		func(value T) error {
			length, isSized := sized.Len(value)

			return verdict.Of(isSized && length == 0, value)
		},
	)
}

// Zero creates a resolution to tell if a value equals the zero value of its
// type: 0, "", false, a nil pointer, or a struct whose fields are all zero.
// Go cannot infer T from the resolution alone, so it is given explicitly:
//
//	be.Zero[int]()
func Zero[T comparable]() screenplay.Resolution[T] {
	return resolution.FromFunc(
		"be zero",
		func(value T) error {
			var zero T

			return verdict.Of(value == zero, value)
		},
	)
}

// True creates a resolution to tell if a boolean is true.
func True() screenplay.Resolution[bool] {
	return resolution.FromFunc(
		"be true",
		func(value bool) error {
			return verdict.Of(value, value)
		},
	)
}

// False creates a resolution to tell if a boolean is false.
func False() screenplay.Resolution[bool] {
	return resolution.FromFunc(
		"be false",
		func(value bool) error {
			return verdict.Of(!value, value)
		},
	)
}
