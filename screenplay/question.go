package screenplay

// Question represents a question that an actor can ask. T is the type of its
// answer, so that asking a question (see Actor.AsksFor) yields a typed value
// without any assertion on the caller's side.
type Question[T any] interface {
	// AnsweredBy returns the answer that an actor provided to the question.
	AnsweredBy(actor *Actor) (T, error)

	// String describes the question.
	String() string
}
