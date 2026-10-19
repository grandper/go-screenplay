package fixture

import "github.com/grandper/go-screenplay/v2/screenplay"

// FakeAbility represents a fake ability that can be given to an actor.
type FakeAbility struct{}

// NewFakeAbility creates a new fake ability.
func NewFakeAbility() *FakeAbility {
	return &FakeAbility{}
}

// Forget cleans up the ability; there is nothing to clean up.
func (*FakeAbility) Forget() error {
	return nil
}

// FakeAbility implements the screenplay.Ability interface.
var _ screenplay.Ability = (*FakeAbility)(nil)
