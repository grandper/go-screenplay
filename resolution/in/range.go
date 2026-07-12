package in

import (
	"cmp"
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"
)

// Range creates a matcher to tell if a number is within a given range.
func Range[T cmp.Ordered](lowerBound, upperBound T) *RangeResolution[T] {
	return &RangeResolution[T]{
		lowerBound: lowerBound,
		upperBound: upperBound,
	}
}

// RangeResolution is a matcher to tell if a number is within a given range.
type RangeResolution[T cmp.Ordered] struct {
	lowerBound T
	upperBound T
}

// Resolve creates a matcher to make an assertion.
func (r *RangeResolution[T]) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		objValue, ok := obj.(T)
		if !ok {
			return false, errors.New("the object should be a comparable type")
		}

		if objValue >= r.lowerBound && objValue <= r.upperBound {
			return true, nil
		}

		return false, nil
	}
}

// String describe the resolution's expectation.
func (r *RangeResolution[T]) String() string {
	return fmt.Sprintf("in the range [%v, %v]", r.lowerBound, r.upperBound)
}

// RangeResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*RangeResolution[int])(nil)
