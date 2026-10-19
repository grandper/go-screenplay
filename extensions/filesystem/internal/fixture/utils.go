package fixture

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

// CreateFilename creates a unique filename in a temporary directory for testing purposes.
func CreateFilename(t *testing.T) string {
	return fmt.Sprintf("%s/foobar-%d.txt", t.TempDir(), uniqueNumber())
}

// CreateDirectoryName creates a unique directory name in a temporary directory for testing purposes.
func CreateDirectoryName(t *testing.T) string {
	return fmt.Sprintf("%s/foobar-%d", t.TempDir(), uniqueNumber())
}

// uniqueNumber draws the number that makes a test file or directory name unique.
func uniqueNumber() uint {
	return rand.Uint() //nolint:gosec // a test file name needs no secure randomness
}
