// Package fixture holds the helpers to test resolutions: a fake resolution
// that every value satisfies or that none does, and the assertions that a
// value satisfies a resolution or not.
//
//	fixture.AssertMatch(t, be.EqualTo(200), 200)
//	fixture.AssertNoMatch(t, be.EqualTo(200), 404)
package fixture
