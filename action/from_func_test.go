package action_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/screenplay/fixture"
)

func TestFromFunc(t *testing.T) {
	adam := screenplay.ActorNamed("Adam")

	doNothing := action.FromFunc("do nothing", func(_ *screenplay.Actor) error {
		return nil
	})
	assert.Implements(t, (*screenplay.Performable)(nil), doNothing)
	assert.Equal(t, "do nothing", doNothing.String())
	require.NoError(t, adam.AttemptsTo(doNothing))

	failToDoSomething := action.FromFunc("fail to do something", func(_ *screenplay.Actor) error {
		return assert.AnError
	})
	assert.Implements(t, (*screenplay.Performable)(nil), failToDoSomething)
	assert.Equal(t, "fail to do something", failToDoSomething.String())
	require.Error(t, adam.AttemptsTo(failToDoSomething))

	t.Run("and questions", func(t *testing.T) {
		t.Run("performs the task with the answers from the questions", func(t *testing.T) {
			urlQuestion := fixture.NewFakeQuestion("URL", "https://example.com")
			portQuestion := fixture.NewFakeQuestion("port", 8080)

			var capturedURL string
			var capturedPort int
			connect := func(_ *screenplay.Actor, url string, port int) screenplay.Performable {
				return action.FromFunc("connect", func(_ *screenplay.Actor) error {
					capturedURL = url
					capturedPort = port
					return nil
				})
			}

			performable := action.FromFuncAndQuestions(
				"connect to %s on port %d",
				connect,
				urlQuestion, portQuestion,
			)
			assert.Implements(t, (*screenplay.Performable)(nil), performable)
			require.NoError(t, adam.AttemptsTo(performable))
			assert.Equal(t, "https://example.com", capturedURL)
			assert.Equal(t, 8080, capturedPort)
		})

		t.Run("implements the stringer interface", func(t *testing.T) {
			performable := action.FromFuncAndQuestions(
				"connect to the server",
				func(_ *screenplay.Actor, url string, port int) screenplay.Performable {
					return action.FromFunc("connect", func(_ *screenplay.Actor) error {
						_, _ = url, port
						return nil
					})
				},
				fixture.NewFakeQuestion("URL", "https://example.com"),
				fixture.NewFakeQuestion("port", 8080),
			)

			assert.Equal(t, "connect to the server", performable.String())
			require.NoError(t, adam.AttemptsTo(performable))
			assert.Equal(t, "connect to the server", performable.String())
		})

		t.Run("converts the answer when it is convertible to the expected type", func(t *testing.T) {
			var capturedPort int64
			performable := action.FromFuncAndQuestions(
				"capture port %d",
				func(_ *screenplay.Actor, port int64) screenplay.Performable {
					return action.FromFunc("capture port", func(_ *screenplay.Actor) error {
						capturedPort = port
						return nil
					})
				},
				fixture.NewFakeQuestion("port", 8080),
			)
			require.NoError(t, adam.AttemptsTo(performable))
			assert.Equal(t, int64(8080), capturedPort)
		})

		t.Run("returns ErrTaskFuncNotFunction when taskFunc is not a function", func(t *testing.T) {
			performable := action.FromFuncAndQuestions(
				"connect to %s",
				"not a function",
				fixture.NewFakeQuestion("URL", "https://example.com"),
			)
			require.ErrorIs(t, adam.AttemptsTo(performable), action.ErrTaskFuncNotFunction)
		})

		t.Run("returns ErrNoQuestions when no questions are provided", func(t *testing.T) {
			performable := action.FromFuncAndQuestions(
				"do nothing useful",
				func(_ *screenplay.Actor) screenplay.Performable {
					return action.FromFunc("doNothing", func(_ *screenplay.Actor) error { return nil })
				},
			)
			require.ErrorIs(t, adam.AttemptsTo(performable), action.ErrNoQuestions)
		})

		t.Run("returns ErrFirstParamNotActor when the first parameter is not *Actor", func(t *testing.T) {
			performable := action.FromFuncAndQuestions(
				"connect to %s",
				func(url string) screenplay.Performable {
					return action.FromFunc("connect", func(_ *screenplay.Actor) error { _ = url; return nil })
				},
				fixture.NewFakeQuestion("URL", "https://example.com"),
			)
			require.ErrorIs(t, adam.AttemptsTo(performable), action.ErrFirstParamNotActor)
		})

		t.Run("returns ErrQuestionCountMismatch when there are too few questions", func(t *testing.T) {
			performable := action.FromFuncAndQuestions(
				"connect to %s on port %d",
				func(_ *screenplay.Actor, url string, port int) screenplay.Performable {
					return action.FromFunc("connect", func(_ *screenplay.Actor) error {
						_, _ = url, port
						return nil
					})
				},
				fixture.NewFakeQuestion("URL", "https://example.com"),
			)
			require.ErrorIs(t, adam.AttemptsTo(performable), action.ErrQuestionCountMismatch)
		})

		t.Run("returns ErrQuestionCountMismatch when there are too many questions", func(t *testing.T) {
			performable := action.FromFuncAndQuestions(
				"connect to %s",
				func(_ *screenplay.Actor, url string) screenplay.Performable {
					return action.FromFunc("connect", func(_ *screenplay.Actor) error { _ = url; return nil })
				},
				fixture.NewFakeQuestion("URL", "https://example.com"),
				fixture.NewFakeQuestion("port", 8080),
			)
			require.ErrorIs(t, adam.AttemptsTo(performable), action.ErrQuestionCountMismatch)
		})

		t.Run("returns ErrTaskFuncMustReturnOne when taskFunc does not return exactly one value", func(t *testing.T) {
			performable := action.FromFuncAndQuestions(
				"connect to %s",
				func(_ *screenplay.Actor, url string) { _ = url },
				fixture.NewFakeQuestion("URL", "https://example.com"),
			)
			require.ErrorIs(t, adam.AttemptsTo(performable), action.ErrTaskFuncMustReturnOne)
		})

		t.Run("returns ErrTaskFuncMustReturnTask when taskFunc does not return a Performable", func(t *testing.T) {
			performable := action.FromFuncAndQuestions(
				"connect to %s",
				func(_ *screenplay.Actor, url string) string { return url },
				fixture.NewFakeQuestion("URL", "https://example.com"),
			)
			require.ErrorIs(t, adam.AttemptsTo(performable), action.ErrTaskFuncMustReturnTask)
		})

		t.Run("returns the underlying error when a question fails to answer", func(t *testing.T) {
			failingQuestion := fixture.NewFailingFakeQuestion[string]("the URL", assert.AnError)
			performable := action.FromFuncAndQuestions(
				"connect to %s",
				func(_ *screenplay.Actor, url string) screenplay.Performable {
					return action.FromFunc("connect", func(_ *screenplay.Actor) error { _ = url; return nil })
				},
				failingQuestion,
			)
			require.ErrorIs(t, adam.AttemptsTo(performable), assert.AnError)
		})

		t.Run(
			"returns ErrAnswerNotAssignable when the answer cannot be assigned to the expected type",
			func(t *testing.T) {
				performable := action.FromFuncAndQuestions(
					"connect to %v",
					func(_ *screenplay.Actor, url string) screenplay.Performable {
						return action.FromFunc("connect", func(_ *screenplay.Actor) error { _ = url; return nil })
					},
					fixture.NewFakeQuestion("URL", struct{}{}),
				)
				require.ErrorIs(t, adam.AttemptsTo(performable), action.ErrAnswerNotAssignable)
			},
		)

		t.Run("returns ErrNotAQuestion when an argument is not a question", func(t *testing.T) {
			performable := action.FromFuncAndQuestions(
				"connect to %s",
				func(_ *screenplay.Actor, url string) screenplay.Performable {
					return action.FromFunc("connect", func(_ *screenplay.Actor) error { _ = url; return nil })
				},
				"not a question",
			)
			require.ErrorIs(t, adam.AttemptsTo(performable), action.ErrNotAQuestion)
		})

		t.Run("returns ErrAnswerNotAssignable when the answer is nil", func(t *testing.T) {
			performable := action.FromFuncAndQuestions(
				"connect to %s",
				func(_ *screenplay.Actor, url string) screenplay.Performable {
					return action.FromFunc("connect", func(_ *screenplay.Actor) error { _ = url; return nil })
				},
				fixture.NewFakeQuestion[any]("the URL", nil),
			)
			require.ErrorIs(t, adam.AttemptsTo(performable), action.ErrAnswerNotAssignable)
		})

		t.Run("returns the error from the inner task", func(t *testing.T) {
			performable := action.FromFuncAndQuestions(
				"connect to %s",
				func(_ *screenplay.Actor, url string) screenplay.Performable {
					return action.FromFunc("connect", func(_ *screenplay.Actor) error {
						_ = url
						return assert.AnError
					})
				},
				fixture.NewFakeQuestion("URL", "https://example.com"),
			)
			require.ErrorIs(t, adam.AttemptsTo(performable), assert.AnError)
		})
	})

	t.Run("fails instead of panicking when the function is nil", func(t *testing.T) {
		err := adam.AttemptsTo(action.FromFunc("open the home page", nil))
		require.EqualError(t, err, "the function given to FromFunc is nil")
	})
}
