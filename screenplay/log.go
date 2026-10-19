package screenplay

import "fmt"

// Log answers a question and narrates its value as an aside, which is handy for
// debugging. Mirroring ScreenPy's Log, it answers the question directly —
// without the "asks for" narration of Sees — and whispers the value at the
// current depth. The aside is rendered by whatever adapters are attached to the
// actor's narrator; with none attached it stays silent.
func Log[T any](question Question[T]) Performable {
	return &logAction[T]{
		question: question,
	}
}

// logAction is an action that narrates the answer to a question.
type logAction[T any] struct {
	question Question[T]
}

// String describes the action.
func (a *logAction[T]) String() string {
	return fmt.Sprintf("log the %s", a.question.String())
}

// PerformAs answers the question and, when the actor is narrating, whispers its
// value as an aside. It returns any error the question produced.
func (a *logAction[T]) PerformAs(actor *Actor) error {
	answer, err := a.question.AnsweredBy(actor)
	if err != nil {
		return err
	}

	if actor.narrates() {
		actor.narrator().WhispersTheAside(actor.name, fmt.Sprintf("=> %v", answer))
	}

	return nil
}

// logAction implements the Performable interface.
var _ Performable = (*logAction[any])(nil)
