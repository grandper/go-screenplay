package resolution

import (
	"errors"

	"github.com/grandper/go-screenplay/v2/internal/verdict"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// FromFunc creates a resolution from a description and the function that
// inspects a value. The function has the signature of
// screenplay.Resolution.Inspect: it returns nil when the value satisfies the
// resolution, and an error explaining why when it does not. The type of the
// values the resolution inspects is inferred from the function.
//
// The description is written to follow "should":
//
//	evenNumber := resolution.FromFunc("be an even number", func(value int) error {
//		if value%2 != 0 {
//			return fmt.Errorf("got %d", value)
//		}
//
//		return nil
//	})
//
// An error returned by the function means that the value does not satisfy the
// resolution. To report that the resolution itself was built wrongly, return
// an error wrapping screenplay.ErrInvalidResolution: every verb then fails,
// the negative ones included. A resolution created with a nil function is
// such an invalid resolution.
func FromFunc[T any](description string, inspect func(value T) error) screenplay.Resolution[T] {
	if inspect == nil {
		return verdict.Invalid[T](description, errors.New("the inspect function is nil"))
	}

	return &funcResolution[T]{
		description: description,
		inspect:     inspect,
	}
}

type funcResolution[T any] struct {
	description string
	inspect     func(value T) error
}

// Inspect returns nil when the value satisfies the resolution.
func (r *funcResolution[T]) Inspect(value T) error {
	return r.inspect(value)
}

// String describes the resolution's expectation.
func (r *funcResolution[T]) String() string {
	return r.description
}

// funcResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution[any] = (*funcResolution[any])(nil)
