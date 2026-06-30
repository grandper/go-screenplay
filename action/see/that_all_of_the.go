package see

import "github.com/grandper/go-screenplay/screenplay"

// ThatAllOfThe creates an action to see if all of the questions match a resolution.
func ThatAllOfThe(questions ...screenplay.Question) func(screenplay.Resolution) screenplay.Performable {
	return func(resolution screenplay.Resolution) screenplay.Performable {
		tuples := make([]any, 0, len(questions))
		for _, question := range questions {
			tuples = append(tuples, question, resolution)
		}

		return AllOf(tuples...)
	}
}
