package see

import (
	"github.com/grandper/go-screenplay/resolution/is"
	"github.com/grandper/go-screenplay/screenplay"
)

// ThatNoneOfThe creates an action to see if none of the questions match a resolution.
func ThatNoneOfThe(questions ...screenplay.Question) func(screenplay.Resolution) screenplay.Performable {
	return func(resolution screenplay.Resolution) screenplay.Performable {
		return ThatAllOfThe(questions...)(is.Not(resolution))
	}
}
