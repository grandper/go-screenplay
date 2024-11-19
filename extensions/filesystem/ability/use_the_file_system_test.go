package ability_test

import (
	"io/fs"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/extensions/filesystem/ability"
)

func TestUseTheFileSystem(t *testing.T) {
	currentFolder, errGetWd := os.Getwd()
	require.NoError(t, errGetWd)

	useTheFileSystem := ability.UseTheFileSystem()

	count := 0
	var mutex sync.Mutex
	getFilename := func() string {
		mutex.Lock()
		defer mutex.Unlock()
		count++
		return "foobar" + strconv.Itoa(count) + ".txt"
	}

	t.Run("should provide the current path", func(t *testing.T) {
		path, err := useTheFileSystem.CurrentPath()
		require.NoError(t, err)
		assert.Equal(t, currentFolder, path)
	})

	t.Run("can change the path", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			path, err := useTheFileSystem.CurrentPath()
			require.NoError(t, err)

			require.NoError(t, useTheFileSystem.ChangeDirectory("."))
			newPath, err := useTheFileSystem.CurrentPath()
			require.NoError(t, err)
			assert.Equal(t, path, newPath)

			require.NoError(t, useTheFileSystem.ChangeDirectory(".."))
			newPath, err = useTheFileSystem.CurrentPath()
			require.NoError(t, err)
			assert.NotEqual(t, path, newPath)

			require.NoError(t, useTheFileSystem.ChangeDirectory("ability"))
			newPath, err = useTheFileSystem.CurrentPath()
			require.NoError(t, err)
			assert.Equal(t, path, newPath)
		})

		t.Run("fails if the path does not exist", func(t *testing.T) {
			require.Error(t, useTheFileSystem.ChangeDirectory("foobar"))
		})
	})

	t.Run("should tell if a file exists", func(t *testing.T) {
		t.Run("the file exists", func(t *testing.T) {
			exists, errExist := useTheFileSystem.FileExists("use_the_file_system_test.go")
			require.NoError(t, errExist)
			assert.True(t, exists)
		})

		t.Run("the file does not exist", func(t *testing.T) {
			exists, errExist := useTheFileSystem.FileExists(getFilename())
			require.NoError(t, errExist)
			assert.False(t, exists)
		})
	})

	t.Run("should tell if a directory exists", func(t *testing.T) {
		t.Run("the directory exists", func(t *testing.T) {
			exists, errExist := useTheFileSystem.DirectoryExists(".")
			require.NoError(t, errExist)
			assert.True(t, exists)
		})

		t.Run("the directory does not exist", func(t *testing.T) {
			exists, errExist := useTheFileSystem.DirectoryExists("foobar")
			require.NoError(t, errExist)
			assert.False(t, exists)
		})
	})

	t.Run("should provide a way to create a temporary file", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			filename, err := useTheFileSystem.CreateTemporaryFile("foo*.txt")
			require.NoError(t, err)

			exists, err := useTheFileSystem.FileExists(filename)
			require.NoError(t, err)
			assert.True(t, exists)
		})

		t.Run("fails if the file exists", func(t *testing.T) {
			filename, err := useTheFileSystem.CreateTemporaryFile("foo*.txt")
			require.NoError(t, err)

			_, err = useTheFileSystem.CreateTemporaryFile(filename)
			require.Error(t, err)
		})
	})

	t.Run("creates a file", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			filename := getFilename()
			err := useTheFileSystem.CreateFile(filename)
			require.NoError(t, err)

			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()

			exists, err := useTheFileSystem.FileExists(filename)
			require.NoError(t, err)
			assert.True(t, exists)
		})

		t.Run("fails if the file cannot be created", func(t *testing.T) {
			err := useTheFileSystem.CreateFile("/foobar.txt")
			require.Error(t, err)
		})

		t.Run("fails if the file already exists", func(t *testing.T) {
			filename := getFilename()
			err := useTheFileSystem.CreateFile(filename)
			require.NoError(t, err)

			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()

			err = useTheFileSystem.CreateFile(filename)
			require.Error(t, err)
		})
	})

	t.Run("creates a file with content", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			filename := getFilename()
			err := useTheFileSystem.CreateFileWithContent(filename, strings.NewReader("hello world"))
			require.NoError(t, err)

			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()

			exists, err := useTheFileSystem.FileExists(filename)
			require.NoError(t, err)
			assert.True(t, exists)

			content, err := useTheFileSystem.Read(filename)
			require.NoError(t, err)
			assert.Equal(t, "hello world", string(content))
		})

		t.Run("fails if the file cannot be created", func(t *testing.T) {
			err := useTheFileSystem.CreateFileWithContent("/foobar.txt", strings.NewReader("hello world"))
			require.Error(t, err)
		})

		t.Run("fails if the content cannot be copied", func(t *testing.T) {
			filename := getFilename()
			defer func() {
				assert.NoError(t, useTheFileSystem.RemoveTheFile(filename))
			}()

			err := useTheFileSystem.CreateFileWithContent(filename, &failingReader{})
			require.Error(t, err)
		})

		t.Run("fails if the file already exists", func(t *testing.T) {
			filename := getFilename()
			err := useTheFileSystem.CreateFile(filename)
			require.NoError(t, err)

			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()

			err = useTheFileSystem.CreateFileWithContent(filename, strings.NewReader("hello world"))
			require.ErrorIs(t, err, ability.ErrFileAlreadyExists)
		})
	})

	t.Run("appends to a file", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			filename := getFilename()
			err := useTheFileSystem.CreateFileWithContent(filename, strings.NewReader("hello"))
			require.NoError(t, err)

			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()

			err = useTheFileSystem.AppendTo(filename, strings.NewReader(" world"))
			require.NoError(t, err)

			content, err := useTheFileSystem.Read(filename)
			require.NoError(t, err)
			assert.Equal(t, "hello world", string(content))
		})

		t.Run("fails if the file does not exist", func(t *testing.T) {
			err := useTheFileSystem.AppendTo(getFilename(), strings.NewReader(" world"))
			require.Error(t, err)
		})

		t.Run("fails if the content cannot be appended", func(t *testing.T) {
			filename := getFilename()
			err := useTheFileSystem.CreateFileWithContent(filename, strings.NewReader("hello"))
			require.NoError(t, err)

			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()

			err = useTheFileSystem.AppendTo(filename, &failingReader{})
			require.Error(t, err)
		})
	})

	t.Run("should remove files", func(t *testing.T) {
		filename, err := useTheFileSystem.CreateTemporaryFile("foo*.bar")
		require.NoError(t, err)

		exists, err := useTheFileSystem.FileExists(filename)
		require.NoError(t, err)
		assert.True(t, exists)

		require.NoError(t, useTheFileSystem.RemoveTheFile(filename))

		exists, err = useTheFileSystem.FileExists(filename)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("create and remove directory", func(t *testing.T) {
		require.NoError(t, useTheFileSystem.CreateDirectory("test", fs.ModePerm))
		require.NoError(t, useTheFileSystem.RemoveDirectory("test"))
	})

	t.Run("should be able to forget", func(t *testing.T) {
		require.NoError(t, useTheFileSystem.Forget())
	})
}

type failingReader struct{}

// Read always fails.
func (r *failingReader) Read(_ []byte) (int, error) {
	return 0, os.ErrInvalid
}
