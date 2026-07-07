package question

import (
	"github.com/grandper/go-screenplay/screenplay"
)

// Silently wraps a question so that the narration produced while it is answered
// is muted: any beats or asides emitted while computing the answer never reach
// the narrator's adapters.
//
// It is meant to keep secrets out of the logs — a question that reads a token or
// a password can be answered without its intermediate steps reaching an adapter
// — and to trim the noise of a chatty question down to nothing.
//
// A production configured with screenplay.WithForceAllNarration neutralises
// this decorator, so a scenario can be debugged with every step narrated again.
func Silently(q screenplay.Question) screenplay.Question {
	return &SilentlyQuestion{
		question: q,
	}
}

// SilentlyQuestion is a question that mutes the narration produced while the
// wrapped question is answered.
type SilentlyQuestion struct {
	question screenplay.Question
}

// String describes the question.
func (q *SilentlyQuestion) String() string {
	return "silently " + q.question.String()
}

// AnsweredBy answers the wrapped question through a muted view of the actor, so
// the narration it would produce is dropped. The view leaves the actor itself
// untouched, so a Silently running in one concurrent branch never mutes
// another. When the actor's narrator cannot be muted (see
// screenplay.WithForceAllNarration), the steps narrate as usual.
func (q *SilentlyQuestion) AnsweredBy(actor *screenplay.Actor) (any, error) {
	return q.question.AnsweredBy(actor.Muted())
}

// SilentlyQuestion implements the screenplay.Question interface.
var _ screenplay.Question = (*SilentlyQuestion)(nil)
