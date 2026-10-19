package ability_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/ability"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action/create"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestUseTheFileSystem(t *testing.T) {
	t.Run("gives an actor the ability to use the file system", func(t *testing.T) {
		adam := screenplay.ActorNamed("Adam").WhoCan(ability.UseTheFileSystem())
		require.NoError(t, adam.AttemptsTo(create.TheDirectory(t.TempDir()+"/orders")))
	})
}
