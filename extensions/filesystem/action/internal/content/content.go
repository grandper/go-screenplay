// Package content holds what the file system actions share about the content
// they write: buffering a reader and describing the content.
package content

import (
	"fmt"
	"io"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/verdict"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// Read buffers the content of a reader once, up to
// screenplay.DefaultMaxBufferedSize, so that an action holding it can be
// performed more than once: a reader is consumed by its first use.
func Read(reader io.Reader) ([]byte, error) {
	content, err := io.ReadAll(io.LimitReader(reader, screenplay.DefaultMaxBufferedSize+1))
	if err != nil {
		return nil, err
	}
	if len(content) > screenplay.DefaultMaxBufferedSize {
		return nil, fmt.Errorf(
			"the content exceeds %d MiB",
			screenplay.DefaultMaxBufferedSize>>20, //nolint:mnd // shifting by 20 converts bytes to MiB
		)
	}
	return content, nil
}

// Describe describes the content of an action as its constructor reads,
// followed by the content itself: the text "Hello", the bytes "Hello". A
// content that could not be read is described without its value.
func Describe(kind string, content []byte, err error) string {
	if err != nil {
		return "the " + kind
	}

	return fmt.Sprintf("the %s %s", kind, verdict.Describe(content))
}
