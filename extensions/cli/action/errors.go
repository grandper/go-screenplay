package action

import "errors"

// ErrEmptyCommandName is returned when the command name is empty.
var ErrEmptyCommandName = errors.New("the command name cannot be empty")
