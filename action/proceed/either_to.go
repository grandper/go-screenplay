package proceed

import (
	impl "github.com/grandper/go-screenplay/v2/action/internal/proceed"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

var (
	// ErrNoPerformableProvided is returned when EitherTo is performed without performables.
	ErrNoPerformableProvided = impl.ErrNoPerformableProvided
	// ErrNoAlternateProvided is returned when EitherTo is performed without alternate performables.
	ErrNoAlternateProvided = impl.ErrNoAlternateProvided
)

// EitherTo performs the provided actions or tasks, or the alternate ones given
// to Or when the first ones fail:
//
//	proceed.EitherTo(RestoreTheSession()).Or(LogInAs("wanda"))
func EitherTo(actions ...screenplay.Performable) *impl.EitherToAction {
	return impl.EitherTo(actions...)
}
