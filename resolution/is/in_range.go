package is

import (
	"cmp"
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"
)

// InRange creates a matcher to tell if a number is within a given range.
func InRange[T cmp.Ordered](lowerBound, upperBound T) *InRangeResolution[T] {
	return &InRangeResolution[T]{
		lowerBound: lowerBound,
		upperBound: upperBound,
	}
}

// InRangeResolution is a matcher to tell if a number is within a given range.
type InRangeResolution[T cmp.Ordered] struct {
	lowerBound T
	upperBound T
}

// Resolve creates a matcher to make an assertion.
func (r *InRangeResolution[T]) Resolve() screenplay.Matcher {
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
func (r *InRangeResolution[T]) String() string {
	return fmt.Sprintf("in the range [%v, %v]", r.lowerBound, r.upperBound)
}

// InRangeResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*InRangeResolution[int])(nil)
