package resolution

import (
	"errors"

	"github.com/grandper/go-screenplay/extensions/filesystem/question"
	"github.com/grandper/go-screenplay/screenplay"
)

var (
	// ErrInvalidObjectType is returned when the object is not a File.
	ErrInvalidObjectType = errors.New("invalid object type: expect a File")
)

// Exists is a matcher to tell if a file exists.
func Exists() *ExistsResolution {
	return &ExistsResolution{}
}

// ExistsResolution is a matcher to tell if a file exists.
type ExistsResolution struct{}

// Resolve creates a matcher to make an assertion.
func (r *ExistsResolution) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		switch value := obj.(type) {
		case *question.File:
			return value.Exists(), nil
		case *question.Directory:
			return value.Exists(), nil
		case nil:
			return false, nil
		default:
			return false, ErrInvalidObjectType
		}
	}
}

// String describe the resolution's expectation.
func (r *ExistsResolution) String() string {
	return "exists"
}
