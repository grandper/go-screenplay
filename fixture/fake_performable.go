package fixture

import "github.com/grandper/go-screenplay/screenplay"

// FakePerformable represents a fake performable that can be performed by an actor.
type FakePerformable struct {
	err         error
	description string
}

// NewFakePerformable creates a new fake performable.
func NewFakePerformable(description string, err error) *FakePerformable {
	return &FakePerformable{
		err:         err,
		description: description,
	}
}

// PerformAs performs the task or the action as the provided actor.
func (fp *FakePerformable) PerformAs(_ *screenplay.Actor) error {
	return fp.err
}

// String describes the action.
func (fp *FakePerformable) String() string {
	return fp.description
}

// FakePerformable implements the screenplay.Performable interface.
var _ screenplay.Performable = (*FakePerformable)(nil)
