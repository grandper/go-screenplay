package contains

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/grandper/go-screenplay/screenplay"
)

// TheBytes creates a matcher to tell if a slice of bytes contains a given slice of bytes.
func TheBytes(b []byte) *TheBytesResolution {
	return &TheBytesResolution{
		bytes: b,
	}
}

// TheBytesResolution is a matcher to tell if a slice of bytes contains a given slice of bytes.
type TheBytesResolution struct {
	bytes []byte
}

// Resolve creates a matcher to make an assertion.
func (r *TheBytesResolution) Resolve() screenplay.Matcher {
	return func(obj any) (bool, error) {
		objValue, ok := obj.([]byte)
		if !ok {
			return false, errors.New("the object should be a slice of bytes")
		}

		if bytes.Contains(objValue, r.bytes) {
			return true, nil
		}

		return false, nil
	}
}

// String describe the resolution's expectation.
func (r *TheBytesResolution) String() string {
	return fmt.Sprintf("containing the bytes %s", r.bytes)
}

// TheBytesResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*TheBytesResolution)(nil)
