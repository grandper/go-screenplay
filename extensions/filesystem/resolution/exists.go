package resolution

import (
	"errors"
	"reflect"

	"github.com/grandper/go-screenplay/v2/resolution"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// Exist creates a resolution to tell if a file or a directory exists. N is the
// type of the value answered by the question, *question.File or
// *question.Directory; Go cannot infer it from the resolution alone, so it is
// given explicitly:
//
//	see.The(FileNamed("notes.txt")).Should(Exist[*question.File]())
//
// A nil node does not exist.
func Exist[N interface{ Exists() bool }]() screenplay.Resolution[N] {
	return resolution.FromFunc(
		"exist",
		func(node N) error {
			if isNil(node) || !node.Exists() {
				return errors.New("it does not exist")
			}

			return nil
		},
	)
}

// isNil reports whether a node is a nil pointer (or a nil interface).
func isNil(node any) bool {
	if node == nil {
		return true
	}

	value := reflect.ValueOf(node)

	return value.Kind() == reflect.Pointer && value.IsNil()
}
