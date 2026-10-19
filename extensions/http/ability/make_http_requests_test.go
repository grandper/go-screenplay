package ability_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/extensions/http/ability"
	"github.com/grandper/go-screenplay/v2/extensions/http/action/send"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestMakeHTTPRequests(t *testing.T) {
	t.Run("gives an actor the ability to make HTTP requests", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		adam := screenplay.ActorNamed("Adam").WhoCan(ability.MakeHTTPRequests().WithTimingOutAfter(5).Seconds())
		require.NoError(t, adam.AttemptsTo(send.GetRequest().To(server.URL+"/orders")))
	})
}
