package screenplay

// Question represents a question that an actor can ask.
type Question interface {
	// AnsweredBy returns the answer that an actor provided to the question.
	AnsweredBy(actor *Actor) (any, error)
	// String describes the question.
	String() string
}
