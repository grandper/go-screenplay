package status_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/extensions/http/ability"
	impl "github.com/grandper/go-screenplay/v2/extensions/http/internal/ability"
	"github.com/grandper/go-screenplay/v2/extensions/http/question/status"
	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/screenplay/fixture"
)

func TestCodeOfQuestion(t *testing.T) {
	t.Parallel()

	adam := screenplay.ActorNamed("Adam")

	httpResponse, setupErr := impl.NewHTTPResponseFrom(&http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader("Hello World")),
	})
	require.NoError(t, setupErr)

	response := fixture.NewFakeQuestion("response", httpResponse)
	nilAnswer := fixture.NewFakeQuestion[*ability.HTTPResponse]("nil answer", nil)
	failingQuestion := fixture.NewFailingFakeQuestion[*ability.HTTPResponse](
		"failing question",
		errors.New("failed to get the answer"),
	)

	t.Run("returns the status code of the HTTP response", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(status.CodeOf(response))

		require.NoError(t, err)
		assert.Equal(t, 200, answer)
	})

	t.Run("fails when the answer is nil", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(status.CodeOf(nilAnswer))

		require.ErrorIs(t, err, status.ErrResponseIsNil)
		assert.Zero(t, answer)
	})

	t.Run("fails when the wrapped question fails", func(t *testing.T) {
		t.Parallel()

		answer, err := adam.AsksFor(status.CodeOf(failingQuestion))

		require.Error(t, err)
		assert.Zero(t, answer)
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		t.Parallel()

		question := status.CodeOf(response)

		assert.Equal(t, "status code of response", question.String())
	})

	t.Run("implements the question interface", func(t *testing.T) {
		t.Parallel()

		assert.Implements(t, (*screenplay.Question[int])(nil), status.CodeOf(response))
	})
}
