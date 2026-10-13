package question_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/action/see"
	"github.com/grandper/go-screenplay/v2/extensions/http/ability"
	"github.com/grandper/go-screenplay/v2/extensions/http/action"
	"github.com/grandper/go-screenplay/v2/extensions/http/question"
	"github.com/grandper/go-screenplay/v2/extensions/http/question/status"
	"github.com/grandper/go-screenplay/v2/question/last"
	"github.com/grandper/go-screenplay/v2/resolution/equal"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestResponsesQuestion(t *testing.T) {
	t.Run("returns all the responses received by the actor", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, err := w.Write([]byte("Hello World"))
			assert.NoError(t, err)
		}))
		defer server.Close()
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(t, adam.AttemptsTo(action.SendGetRequest().To(server.URL)))
		assert.NoError(t, adam.AttemptsTo(action.SendGetRequest().To(server.URL)))

		responses, err := adam.AsksFor(question.Responses())

		require.NoError(t, err)
		assert.Len(t, responses, 2)
	})

	t.Run("can be composed to ask about the status code of the last response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, err := w.Write([]byte("Hello World"))
			assert.NoError(t, err)
		}))
		defer server.Close()
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(t, adam.AttemptsTo(action.SendGetRequest().To(server.URL)))
		assert.NoError(
			t,
			adam.AttemptsTo(see.The(status.CodeOf(last.Of(question.Responses()))).Is(equal.To(200))),
		)
	})

	t.Run("fails when the actor does not have the ability MakeHTTPRequests", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")

		_, err := adam.AsksFor(question.Responses())

		assert.Error(t, err)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "HTTP responses", question.Responses().String())
	})

	t.Run("implements the question interface", func(t *testing.T) {
		assert.Implements(t, (*screenplay.Question[[]*ability.HTTPResponse])(nil), question.Responses())
	})
}
