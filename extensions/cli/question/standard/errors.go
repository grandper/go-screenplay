package standard

import "errors"

// ErrResultIsNil is returned when the wrapped question answers with a nil CLI result.
var ErrResultIsNil = errors.New("the answer of the question is a nil CLI result")
