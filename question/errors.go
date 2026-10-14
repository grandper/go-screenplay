package question

import "errors"

var (
	// ErrAnswerIsEmpty is returned by first.Of and last.Of when the wrapped
	// question answers with an empty slice.
	ErrAnswerIsEmpty = errors.New("the answer of the question is an empty slice")

	// ErrAnswerIsNotAStruct is returned by body.Of, data.Of and header.Of when
	// the wrapped question does not answer with a struct.
	ErrAnswerIsNotAStruct = errors.New("the answer of the question is not a struct")

	// ErrFieldNotFound is returned by body.Of, data.Of and header.Of when the
	// struct answered by the wrapped question has no such field.
	ErrFieldNotFound = errors.New("the answer of the question has no such field")

	// ErrFieldHasWrongType is returned by body.Of, data.Of and header.Of when
	// the field of the struct answered by the wrapped question is not of the
	// expected type.
	ErrFieldHasWrongType = errors.New("the field of the answer of the question has the wrong type")
)
