package question_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/action/see"
	"github.com/grandper/go-screenplay/extensions/http/ability"
	"github.com/grandper/go-screenplay/extensions/http/action"
	"github.com/grandper/go-screenplay/extensions/http/question"
	"github.com/grandper/go-screenplay/extensions/http/question/status"
	"github.com/grandper/go-screenplay/question/last"
	"github.com/grandper/go-screenplay/resolution/is"
	"github.com/grandper/go-screenplay/screenplay"
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

		answer, err := adam.AsksFor(question.Responses())

		require.NoError(t, err)
		responses, ok := answer.([]*ability.HTTPResponse)
		require.True(t, ok)
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
			adam.AttemptsTo(see.The(status.CodeOf(last.Of(question.Responses()))).Is(is.EqualTo(200))),
		)
	})

	t.Run("fails when the actor does not have the ability MakeHTTPRequests", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam")

		_, err := adam.AsksFor(question.Responses())

		assert.Error(t, err)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		assert.Equal(t, "the HTTP responses", question.Responses().String())
	})
}
