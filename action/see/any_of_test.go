package see_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/action/see"
	"github.com/grandper/go-screenplay/fixture"
	"github.com/grandper/go-screenplay/resolution/contains"
	"github.com/grandper/go-screenplay/resolution/is"
	"github.com/grandper/go-screenplay/resolution/testdata"
	"github.com/grandper/go-screenplay/screenplay"
)

func TestSeeAnyOfAction(t *testing.T) {
	t.Parallel()

	adam := screenplay.ActorNamed("Adam")
	profilWidget := fixture.NewFakeQuestion("profil widget", "<profil>Adam</profil>")
	loginForm := fixture.NewFakeQuestion("login form", "<login>adam@google.com</login>")

	missingProfileWidget := fixture.NewFailingFakeQuestion(
		"profil widget",
		errors.New("failed to get the profil widget"),
	)
	containsTheTextButFails := testdata.NewFailingResolution(
		"contains the text Adam",
		errors.New("failed to match the content of the profil widget"),
	)

	t.Run("should see any of the list items", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, adam.AttemptsTo(see.AnyOf(
			profilWidget, contains.TheText("Adam"),
			loginForm, contains.TheText("adam@google.com"))))

		require.NoError(t, adam.AttemptsTo(see.AnyOf(
			profilWidget, contains.TheText("Foobar"),
			loginForm, contains.TheText("adam@google.com"))))

		require.NoError(t, adam.AttemptsTo(see.AnyOf(
			profilWidget, contains.TheText("Adam"),
			loginForm, contains.TheText("foo@bar.com"))))
	})

	t.Run("will succeed without error when no expectation is provided", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, adam.AttemptsTo(see.AnyOf()))
	})

	t.Run("will succeed even if a few question failed to get answer", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, adam.AttemptsTo(see.AnyOf(
			missingProfileWidget, contains.TheText("Adam"),
			loginForm, contains.TheText("adam@google.com"))))
	})

	t.Run("will succeed even if a few resolution failed to match", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, adam.AttemptsTo(see.AnyOf(
			profilWidget, contains.TheText("Adam"),
			loginForm, containsTheTextButFails)))
	})

	t.Run("fails when there is nothing to see", func(t *testing.T) {
		t.Parallel()
		require.Error(t, adam.AttemptsTo(see.AnyOf(
			profilWidget, contains.TheText("Foobar"),
			loginForm, contains.TheText("foo@bar.com"))))
	})

	t.Run("fails when the argument list is incomplete", func(t *testing.T) {
		t.Parallel()
		require.Error(t, adam.AttemptsTo(see.AnyOf(profilWidget)))

		require.Error(t, adam.AttemptsTo(see.AnyOf(
			profilWidget, contains.TheText("Foobar"),
			loginForm,
		)))
	})

	t.Run("fails when a resolution is passed instead of a question", func(t *testing.T) {
		t.Parallel()
		require.Error(t, adam.AttemptsTo(see.AnyOf(profilWidget, loginForm)))
	})

	t.Run("fails when a question is passed instead of a resolution", func(t *testing.T) {
		t.Parallel()
		require.Error(t, adam.AttemptsTo(see.AnyOf(contains.TheText("Foobar"), profilWidget)))
	})

	t.Run("fails when the actor fails to answer all the questions", func(t *testing.T) {
		t.Parallel()
		require.Error(t, adam.AttemptsTo(see.AnyOf(missingProfileWidget, is.EqualTo("hello everybody"))))
	})

	t.Run("fails when the resolution fails all the resolution", func(t *testing.T) {
		t.Parallel()
		require.Error(t, adam.AttemptsTo(see.AnyOf(profilWidget, containsTheTextButFails)))
	})

	t.Run("implements the stringer interface", func(t *testing.T) {
		t.Parallel()

		action := see.AnyOf(
			profilWidget, contains.TheText("Adam"),
			loginForm, contains.TheText("adam@google.com"))

		assert.Equal(
			t,
			"see if the profil widget is containing the text Adam, or see if the login form is containing the text adam@google.com",
			action.String(),
		)
	})
}
