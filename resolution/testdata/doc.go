// Package testdata holds the helpers to test resolutions: a fake resolution
// that every value satisfies or that none does, and the assertions that a
// value satisfies a resolution or not.
//
//	testdata.AssertMatch(t, be.EqualTo(200), 200)
//	testdata.AssertNoMatch(t, be.EqualTo(200), 404)
package testdata
