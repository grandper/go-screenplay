package question_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/extensions/filesystem/question"
)

func TestCreateFakeExistingFile(t *testing.T) {
	t.Run("should create a fake existing file", func(t *testing.T) {
		filePath := "test.txt"
		file := question.CreateFakeExistingFile(filePath)
		require.NotNil(t, file)
		require.Equal(t, ".", file.Dir())
		require.Equal(t, "test.txt", file.Name())
		require.True(t, file.Exists())
	})
}

func TestCreateFakeMissingFile(t *testing.T) {
	t.Run("should create a fake missing file", func(t *testing.T) {
		filePath := "test.txt"
		file := question.CreateFakeMissingFile(filePath)
		require.NotNil(t, file)
		require.Equal(t, ".", file.Dir())
		require.Equal(t, "test.txt", file.Name())
		require.False(t, file.Exists())
	})
}

func TestCreateFakeExistingDirectory(t *testing.T) {
	t.Run("should create a fake existing directory", func(t *testing.T) {
		directoryPath := "path/to/MyDirectory"
		directory := question.CreateFakeExistingDirectory(directoryPath)
		require.NotNil(t, directory)
		require.Equal(t, "path/to", directory.Dir())
		require.Equal(t, "MyDirectory", directory.Name())
		require.True(t, directory.Exists())
	})
}

func TestCreateFakeMissingDirectory(t *testing.T) {
	t.Run("should create a fake missing directory", func(t *testing.T) {
		directoryPath := "path/to/MyDirectory"
		directory := question.CreateFakeMissingDirectory(directoryPath)
		require.NotNil(t, directory)
		require.Equal(t, "path/to", directory.Dir())
		require.Equal(t, "MyDirectory", directory.Name())
		require.False(t, directory.Exists())
	})
}
