package testdata

import (
	"github.com/grandper/go-screenplay/screenplay"
)

// FakeResolution is a fake matcher for test purpose.
type FakeResolution struct {
	result      bool
	err         error
	description string
}

// NewFakeResolution creates a new fake matcher.
func NewFakeResolution(description string, result bool) *FakeResolution {
	return &FakeResolution{
		result:      result,
		err:         nil,
		description: description,
	}
}

// NewFailingResolution creates a new fake matcher that fails.
func NewFailingResolution(description string, err error) *FakeResolution {
	return &FakeResolution{
		result:      false,
		err:         err,
		description: description,
	}
}

// Resolve creates a matcher to make an assertion.
func (fr *FakeResolution) Resolve() screenplay.Matcher {
	return func(_ any) (bool, error) {
		if fr.err != nil {
			return false, fr.err
		}

		return fr.result, nil
	}
}

// String describe the resolution's expectation.
func (fr *FakeResolution) String() string {
	return fr.description
}

// FakeResolution implements the screenplay.Resolution interface.
var _ screenplay.Resolution = (*FakeResolution)(nil)
