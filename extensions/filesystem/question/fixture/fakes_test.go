package fixture_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/question/fixture"
)

func TestFakes(t *testing.T) {
	t.Run("create fake existing file", func(t *testing.T) {
		t.Run("should create a fake existing file", func(t *testing.T) {
			filePath := "test.txt"
			file := fixture.CreateFakeExistingFile(filePath)
			require.NotNil(t, file)
			require.Equal(t, ".", file.Dir())
			require.Equal(t, "test.txt", file.Name())
			require.True(t, file.Exists())
		})
	})

	t.Run("create fake missing file", func(t *testing.T) {
		t.Run("should create a fake missing file", func(t *testing.T) {
			filePath := "test.txt"
			file := fixture.CreateFakeMissingFile(filePath)
			require.NotNil(t, file)
			require.Equal(t, ".", file.Dir())
			require.Equal(t, "test.txt", file.Name())
			require.False(t, file.Exists())
		})
	})

	t.Run("create fake existing directory", func(t *testing.T) {
		t.Run("should create a fake existing directory", func(t *testing.T) {
			directoryPath := "path/to/MyDirectory"
			directory := fixture.CreateFakeExistingDirectory(directoryPath)
			require.NotNil(t, directory)
			require.Equal(t, "path/to", directory.Dir())
			require.Equal(t, "MyDirectory", directory.Name())
			require.True(t, directory.Exists())
		})
	})

	t.Run("create fake missing directory", func(t *testing.T) {
		t.Run("should create a fake missing directory", func(t *testing.T) {
			directoryPath := "path/to/MyDirectory"
			directory := fixture.CreateFakeMissingDirectory(directoryPath)
			require.NotNil(t, directory)
			require.Equal(t, "path/to", directory.Dir())
			require.Equal(t, "MyDirectory", directory.Name())
			require.False(t, directory.Exists())
		})
	})
}
