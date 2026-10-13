package see

import "github.com/grandper/go-screenplay/v2/screenplay"

// ThatAnyOfThe creates an action to see if any of the questions match a resolution.
func ThatAnyOfThe(questions ...screenplay.Question) func(screenplay.Resolution) screenplay.Performable {
	return func(resolution screenplay.Resolution) screenplay.Performable {
		tuples := make([]any, 0, len(questions))
		for _, question := range questions {
			tuples = append(tuples, question, resolution)
		}

		return AnyOf(tuples...)
	}
}
