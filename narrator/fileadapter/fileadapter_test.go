package fileadapter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/narrator/fileadapter"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

func TestFileAdapter(t *testing.T) {
	t.Parallel()

	beat := screenplay.Event{
		Kind:    screenplay.KindBeat,
		Phase:   screenplay.PhaseBegin,
		Depth:   1,
		Message: "Adam sends a GET request to /orders/42",
	}

	t.Run("writes the indented tree to the file", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(t.TempDir(), "narration.log")
		adapter, err := fileadapter.New(path)
		require.NoError(t, err)

		adapter.Narrate(beat)
		require.NoError(t, adapter.Close())

		content, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, "    Adam sends a GET request to /orders/42\n", string(content))
	})

	t.Run("overwrites the file by default", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(t.TempDir(), "narration.log")
		require.NoError(t, os.WriteFile(path, []byte("previous run\n"), 0o644))

		adapter, err := fileadapter.New(path)
		require.NoError(t, err)

		adapter.Narrate(beat)
		require.NoError(t, adapter.Close())

		content, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, "    Adam sends a GET request to /orders/42\n", string(content))
	})

	t.Run("overwrites the file with WithOverwriting", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(t.TempDir(), "narration.log")
		require.NoError(t, os.WriteFile(path, []byte("previous run\n"), 0o644))

		adapter, err := fileadapter.New(path, fileadapter.WithOverwriting())
		require.NoError(t, err)

		adapter.Narrate(beat)
		require.NoError(t, adapter.Close())

		content, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, "    Adam sends a GET request to /orders/42\n", string(content))
	})

	t.Run("appends to the file with WithAppending", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(t.TempDir(), "narration.log")
		require.NoError(t, os.WriteFile(path, []byte("previous run\n"), 0o644))

		adapter, err := fileadapter.New(path, fileadapter.WithAppending())
		require.NoError(t, err)

		adapter.Narrate(beat)
		require.NoError(t, adapter.Close())

		content, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, "previous run\n    Adam sends a GET request to /orders/42\n", string(content))
	})

	t.Run("prefixes each line with the date and time with WithTimestamps", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(t.TempDir(), "narration.log")
		adapter, err := fileadapter.New(path, fileadapter.WithTimestamps())
		require.NoError(t, err)

		adapter.Narrate(beat)
		require.NoError(t, adapter.Close())

		content, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Regexp(t,
			`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}     Adam sends a GET request to /orders/42\n$`,
			string(content))
	})

	t.Run("fails when both WithOverwriting and WithAppending are given", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(t.TempDir(), "narration.log")
		adapter, err := fileadapter.New(path, fileadapter.WithOverwriting(), fileadapter.WithAppending())

		assert.Nil(t, adapter)
		require.ErrorIs(t, err, fileadapter.ErrConflictingOptions)
		assert.Equal(t, "conflicting options: overwriting and appending to "+path, err.Error())
		assert.NoFileExists(t, path)
	})

	t.Run("fails when the file cannot be opened", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(t.TempDir(), "logs", "narration.log")
		adapter, err := fileadapter.New(path)

		assert.Nil(t, adapter)
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("reports a write error when it is closed", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(t.TempDir(), "narration.log")
		adapter, err := fileadapter.New(path)
		require.NoError(t, err)

		require.NoError(t, adapter.Close())
		adapter.Narrate(beat)

		require.ErrorIs(t, adapter.Close(), os.ErrClosed)
	})
}
