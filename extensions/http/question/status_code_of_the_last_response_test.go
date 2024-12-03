package question_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/action/see"
	"github.com/grandper/go-screenplay/extensions/http/ability"
	"github.com/grandper/go-screenplay/extensions/http/action"
	"github.com/grandper/go-screenplay/extensions/http/question"
	"github.com/grandper/go-screenplay/resolution/is"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestStatusCodeOfTheLastResponseQuestion(t *testing.T) {
	t.Run("returns the headers of the last response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, err := w.Write([]byte("Hello World"))
			assert.NoError(t, err)
		}))
		defer server.Close()
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.NoError(t, adam.AttemptsTo(action.SendGetRequest().To(server.URL)))
		assert.NoError(t, adam.AttemptsTo(see.The(question.StatusCodeOfTheLastResponse(), is.EqualTo(200))))
	})

	t.Run("fails to get the status code of the last response when no request has been made", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests())
		assert.Error(t, adam.AttemptsTo(see.The(question.StatusCodeOfTheLastResponse(), is.EqualTo(200))))
	})

	t.Run(
		"fails to get the status code of the last response if the actor does not have the ability MakeHttpRequest",
		func(t *testing.T) {
			adam := screenplay.ActorNamed("Adam")
			assert.Error(t, adam.AttemptsTo(see.The(question.StatusCodeOfTheLastResponse(), is.EqualTo(200))))
		},
	)

	t.Run("implements the stringer interface", func(t *testing.T) {
		statusCode := question.StatusCodeOfTheLastResponse()
		assert.Equal(t, "HTTP status code of the last response", statusCode.String())
	})
}
