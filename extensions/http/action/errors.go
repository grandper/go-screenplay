package action

import "errors"

// ErrInvalidHeaders is wrapped by the error add.TheHeader, add.TheHeaders,
// set.TheHeader and set.TheHeaders return when they are given headers they
// cannot use: an empty key, no header at all, or a key without its value.
// Nothing panics when the action is built; the mistake is reported when the
// actor performs it.
var ErrInvalidHeaders = errors.New("invalid headers")
