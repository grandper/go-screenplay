package see

import "github.com/grandper/go-screenplay/v2/screenplay"

// ThatNoneOfThe creates an action to see if none of the questions match a resolution.
func ThatNoneOfThe(questions ...screenplay.Question) func(screenplay.Resolution) screenplay.Performable {
	return func(resolution screenplay.Resolution) screenplay.Performable {
		tests := make([]*TheAction, 0, len(questions))
		for _, question := range questions {
			tests = append(tests, The(question).IsNot(resolution))
		}

		return &AllOfAction{tests: tests}
	}
}
