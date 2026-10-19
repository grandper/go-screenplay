package ability_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
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
		path, err := useTheFileSystem.ToRetrieveCurrentPath()
		require.NoError(t, err)
		assert.Equal(t, currentFolder, path)
	})

	t.Run("can change the path", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			path, err := useTheFileSystem.ToRetrieveCurrentPath()
			require.NoError(t, err)

			require.NoError(t, useTheFileSystem.ToChangeDirectory("."))
			newPath, err := useTheFileSystem.ToRetrieveCurrentPath()
			require.NoError(t, err)
			assert.Equal(t, path, newPath)

			require.NoError(t, useTheFileSystem.ToChangeDirectory(".."))
			newPath, err = useTheFileSystem.ToRetrieveCurrentPath()
			require.NoError(t, err)
			assert.NotEqual(t, path, newPath)

			require.NoError(t, useTheFileSystem.ToChangeDirectory("ability"))
			newPath, err = useTheFileSystem.ToRetrieveCurrentPath()
			require.NoError(t, err)
			assert.Equal(t, path, newPath)
		})

		t.Run("fails if the path does not exist", func(t *testing.T) {
			require.Error(t, useTheFileSystem.ToChangeDirectory("foobar"))
		})

		t.Run("fails if the path is not a directory", func(t *testing.T) {
			require.Error(t, useTheFileSystem.ToChangeDirectory("use_the_file_system.go"))
		})

		t.Run("without changing the working directory of the process", func(t *testing.T) {
			other := ability.UseTheFileSystem()
			require.NoError(t, other.ToChangeDirectory(t.TempDir()))

			processDir, err := os.Getwd()
			require.NoError(t, err)
			assert.Equal(t, currentFolder, processDir)

			// Another ability is not affected either.
			path, err := useTheFileSystem.ToRetrieveCurrentPath()
			require.NoError(t, err)
			assert.Equal(t, currentFolder, path)
		})

		t.Run("so that relative paths resolve against the new directory", func(t *testing.T) {
			dir := t.TempDir()
			other := ability.UseTheFileSystem()
			require.NoError(t, other.ToChangeDirectory(dir))

			require.NoError(t, other.ToCreateDirectory("docs", fs.ModePerm))
			require.NoError(t, other.ToCreateFileWithContent("docs/notes.txt", strings.NewReader("hello")))
			require.NoError(t, other.ToAppendTo("docs/notes.txt", strings.NewReader(" world")))

			exists, err := other.ToTellIfTheFileExists("docs/notes.txt")
			require.NoError(t, err)
			assert.True(t, exists)
			content, err := other.ToRead("docs/notes.txt")
			require.NoError(t, err)
			assert.Equal(t, "hello world", string(content))
			assert.FileExists(t, filepath.Join(dir, "docs", "notes.txt"))

			require.NoError(t, other.ToOverwriteFileWithContent("docs/notes.txt", strings.NewReader("bye")))
			require.NoError(t, other.ToRemoveFile("docs/notes.txt"))
			require.NoError(t, other.ToRemoveDirectory("docs"))
			assert.NoDirExists(t, filepath.Join(dir, "docs"))

			// Removing the ability's own current directory is refused.
			require.ErrorIs(t, other.ToRemoveDirectory("."), ability.ErrProtectedDirectory)
		})
	})

	t.Run("should tell if a file exists", func(t *testing.T) {
		t.Run("the file exists", func(t *testing.T) {
			exists, errExist := useTheFileSystem.ToTellIfTheFileExists("use_the_file_system_test.go")
			require.NoError(t, errExist)
			assert.True(t, exists)
		})

		t.Run("the file does not exist", func(t *testing.T) {
			exists, errExist := useTheFileSystem.ToTellIfTheFileExists(getFilename())
			require.NoError(t, errExist)
			assert.False(t, exists)
		})
	})

	t.Run("should tell if a directory exists", func(t *testing.T) {
		t.Run("the directory exists", func(t *testing.T) {
			exists, errExist := useTheFileSystem.ToTellIfTheDirectoryExists(".")
			require.NoError(t, errExist)
			assert.True(t, exists)
		})

		t.Run("the directory does not exist", func(t *testing.T) {
			exists, errExist := useTheFileSystem.ToTellIfTheDirectoryExists("foobar")
			require.NoError(t, errExist)
			assert.False(t, exists)
		})
	})

	t.Run("should provide a way to create a temporary file", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			filename, err := useTheFileSystem.ToCreateTemporaryFile("foo*.txt")
			require.NoError(t, err)

			exists, err := useTheFileSystem.ToTellIfTheFileExists(filename)
			require.NoError(t, err)
			assert.True(t, exists)
		})

		t.Run("fails if the file exists", func(t *testing.T) {
			filename, err := useTheFileSystem.ToCreateTemporaryFile("foo*.txt")
			require.NoError(t, err)

			_, err = useTheFileSystem.ToCreateTemporaryFile(filename)
			require.Error(t, err)
		})
	})

	t.Run("creates a temporary file with content", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			filename, err := useTheFileSystem.ToCreateTemporaryFileWithContent(
				"foo*.txt",
				strings.NewReader("hello world"),
			)
			require.NoError(t, err)

			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()

			content, err := useTheFileSystem.ToRead(filename)
			require.NoError(t, err)
			assert.Equal(t, "hello world", string(content))
		})

		t.Run("in the temporary directory, whatever the files of the current directory", func(t *testing.T) {
			// A file of the same name in the current directory is unrelated: the
			// temporary file gets a random suffix in the temporary directory.
			filename := getFilename()
			require.NoError(t, useTheFileSystem.ToCreateFile(filename))
			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()

			temporary, err := useTheFileSystem.ToCreateTemporaryFileWithContent(
				filename,
				strings.NewReader("hello world"),
			)
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, os.Remove(temporary))
			}()

			assert.NotEqual(t, filename, temporary)
			systemTemporaryDirectory := os.TempDir() //nolint:usetesting // the ability writes there
			assert.Equal(t, filepath.Clean(systemTemporaryDirectory), filepath.Dir(temporary))
		})

		t.Run("fails if the content cannot be copied", func(t *testing.T) {
			filename, err := useTheFileSystem.ToCreateTemporaryFileWithContent("foo*.txt", &failingReader{})
			require.Error(t, err)
			assert.Empty(t, filename)
		})
	})

	t.Run("overwrites a file with content", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			filename := getFilename()
			err := useTheFileSystem.ToCreateFileWithContent(filename, strings.NewReader("hello world"))
			require.NoError(t, err)

			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()

			err = useTheFileSystem.ToOverwriteFileWithContent(filename, strings.NewReader("bye"))
			require.NoError(t, err)

			content, err := useTheFileSystem.ToRead(filename)
			require.NoError(t, err)
			assert.Equal(t, "bye", string(content))
		})

		t.Run("fails if the file does not exist", func(t *testing.T) {
			err := useTheFileSystem.ToOverwriteFileWithContent(getFilename(), strings.NewReader("hello world"))
			require.Error(t, err)
		})

		t.Run("fails if the content cannot be copied", func(t *testing.T) {
			filename := getFilename()
			err := useTheFileSystem.ToCreateFileWithContent(filename, strings.NewReader("hello world"))
			require.NoError(t, err)

			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()

			err = useTheFileSystem.ToOverwriteFileWithContent(filename, &failingReader{})
			require.Error(t, err)
		})
	})

	t.Run("creates a temporary directory", func(t *testing.T) {
		dir, err := useTheFileSystem.ToCreateTemporaryDirectory("foobar")
		require.NoError(t, err)

		defer func() {
			assert.NoError(t, useTheFileSystem.ToRemoveDirectory(dir))
		}()

		exists, err := useTheFileSystem.ToTellIfTheDirectoryExists(dir)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("creates a file", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			filename := getFilename()
			err := useTheFileSystem.ToCreateFile(filename)
			require.NoError(t, err)

			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()

			exists, err := useTheFileSystem.ToTellIfTheFileExists(filename)
			require.NoError(t, err)
			assert.True(t, exists)
		})

		t.Run("fails if the file cannot be created", func(t *testing.T) {
			err := useTheFileSystem.ToCreateFile("/foobar.txt")
			require.Error(t, err)
		})

		t.Run("fails if the file already exists", func(t *testing.T) {
			filename := getFilename()
			err := useTheFileSystem.ToCreateFile(filename)
			require.NoError(t, err)

			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()

			err = useTheFileSystem.ToCreateFile(filename)
			require.ErrorIs(t, err, ability.ErrFileAlreadyExists)
		})

		t.Run("readable and writable by its owner only", func(t *testing.T) {
			filename := getFilename()
			require.NoError(t, useTheFileSystem.ToCreateFile(filename))

			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()

			info, err := os.Stat(filename)
			require.NoError(t, err)
			assert.Equal(t, fs.FileMode(0600), info.Mode().Perm())
		})
	})

	t.Run("never writes through a symbolic link", func(t *testing.T) {
		// A dangling link: the target does not exist, so a naive existence check
		// would see "no file" and a naive write would create the target.
		newDanglingLink := func(t *testing.T) (string, string) {
			t.Helper()
			target := getFilename()
			link := getFilename()
			require.NoError(t, os.Symlink(target, link))
			t.Cleanup(func() {
				assert.NoError(t, os.Remove(link))
				_, err := os.Lstat(target)
				assert.ErrorIs(t, err, os.ErrNotExist, "the target of the link must not have been created")
			})
			return link, target
		}

		t.Run("when creating a file", func(t *testing.T) {
			link, _ := newDanglingLink(t)
			err := useTheFileSystem.ToCreateFile(link)
			require.ErrorIs(t, err, ability.ErrFileAlreadyExists)
		})

		t.Run("when creating a file with content", func(t *testing.T) {
			link, _ := newDanglingLink(t)
			err := useTheFileSystem.ToCreateFileWithContent(link, strings.NewReader("owned"))
			require.ErrorIs(t, err, ability.ErrFileAlreadyExists)
		})

		t.Run("when overwriting a file", func(t *testing.T) {
			link, _ := newDanglingLink(t)
			err := useTheFileSystem.ToOverwriteFileWithContent(link, strings.NewReader("owned"))
			require.ErrorIs(t, err, ability.ErrSymbolicLink)
		})

		t.Run("when appending to a file", func(t *testing.T) {
			link, _ := newDanglingLink(t)
			err := useTheFileSystem.ToAppendTo(link, strings.NewReader("owned"))
			require.ErrorIs(t, err, ability.ErrSymbolicLink)
		})

		t.Run("even when the link points to an existing file", func(t *testing.T) {
			target := getFilename()
			require.NoError(t, useTheFileSystem.ToCreateFileWithContent(target, strings.NewReader("original")))
			link := getFilename()
			require.NoError(t, os.Symlink(target, link))
			t.Cleanup(func() {
				assert.NoError(t, os.Remove(link))
				assert.NoError(t, os.Remove(target))
			})

			require.ErrorIs(
				t,
				useTheFileSystem.ToOverwriteFileWithContent(link, strings.NewReader("owned")),
				ability.ErrSymbolicLink,
			)
			require.ErrorIs(t, useTheFileSystem.ToAppendTo(link, strings.NewReader("owned")), ability.ErrSymbolicLink)

			content, err := useTheFileSystem.ToRead(target)
			require.NoError(t, err)
			assert.Equal(t, "original", string(content))
		})
	})

	t.Run("creates a file with content", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			filename := getFilename()
			err := useTheFileSystem.ToCreateFileWithContent(filename, strings.NewReader("hello world"))
			require.NoError(t, err)

			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()

			exists, err := useTheFileSystem.ToTellIfTheFileExists(filename)
			require.NoError(t, err)
			assert.True(t, exists)

			content, err := useTheFileSystem.ToRead(filename)
			require.NoError(t, err)
			assert.Equal(t, "hello world", string(content))
		})

		t.Run("fails if the file cannot be created", func(t *testing.T) {
			err := useTheFileSystem.ToCreateFileWithContent("/foobar.txt", strings.NewReader("hello world"))
			require.Error(t, err)
		})

		t.Run("fails if the content cannot be copied", func(t *testing.T) {
			filename := getFilename()
			defer func() {
				assert.NoError(t, useTheFileSystem.ToRemoveFile(filename))
			}()

			err := useTheFileSystem.ToCreateFileWithContent(filename, &failingReader{})
			require.Error(t, err)
		})

		t.Run("fails if the file already exists", func(t *testing.T) {
			filename := getFilename()
			err := useTheFileSystem.ToCreateFile(filename)
			require.NoError(t, err)

			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()

			err = useTheFileSystem.ToCreateFileWithContent(filename, strings.NewReader("hello world"))
			require.ErrorIs(t, err, ability.ErrFileAlreadyExists)
		})
	})

	t.Run("appends to a file", func(t *testing.T) {
		t.Run("successfully", func(t *testing.T) {
			filename := getFilename()
			err := useTheFileSystem.ToCreateFileWithContent(filename, strings.NewReader("hello"))
			require.NoError(t, err)

			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()

			err = useTheFileSystem.ToAppendTo(filename, strings.NewReader(" world"))
			require.NoError(t, err)

			content, err := useTheFileSystem.ToRead(filename)
			require.NoError(t, err)
			assert.Equal(t, "hello world", string(content))
		})

		t.Run("fails if the file does not exist", func(t *testing.T) {
			err := useTheFileSystem.ToAppendTo(getFilename(), strings.NewReader(" world"))
			require.Error(t, err)
		})

		t.Run("fails if the content cannot be appended", func(t *testing.T) {
			filename := getFilename()
			err := useTheFileSystem.ToCreateFileWithContent(filename, strings.NewReader("hello"))
			require.NoError(t, err)

			defer func() {
				assert.NoError(t, os.Remove(filename))
			}()

			err = useTheFileSystem.ToAppendTo(filename, &failingReader{})
			require.Error(t, err)
		})
	})

	t.Run("should remove files", func(t *testing.T) {
		filename, err := useTheFileSystem.ToCreateTemporaryFile("foo*.bar")
		require.NoError(t, err)

		exists, err := useTheFileSystem.ToTellIfTheFileExists(filename)
		require.NoError(t, err)
		assert.True(t, exists)

		require.NoError(t, useTheFileSystem.ToRemoveFile(filename))

		exists, err = useTheFileSystem.ToTellIfTheFileExists(filename)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("create and remove directory", func(t *testing.T) {
		require.NoError(t, useTheFileSystem.ToCreateDirectory("test", fs.ModePerm))
		require.NoError(t, useTheFileSystem.ToRemoveDirectory("test"))
	})

	t.Run("refuses to remove a protected directory", func(t *testing.T) {
		home, err := os.UserHomeDir()
		require.NoError(t, err)
		root := currentFolder
		for filepath.Dir(root) != root {
			root = filepath.Dir(root)
		}

		protected := map[string]string{
			"the root of the file system":            root,
			"the home directory":                     home,
			"the current working directory":          ".",
			"the current working directory, by path": currentFolder,
			"the parent of the current directory":    "..",
			"an ancestor of the current directory":   filepath.Dir(filepath.Dir(currentFolder)),
		}
		for name, dir := range protected {
			t.Run(name, func(t *testing.T) {
				errRemove := useTheFileSystem.ToRemoveDirectory(dir)
				require.ErrorIs(t, errRemove, ability.ErrProtectedDirectory)

				_, errStat := os.Stat(dir)
				require.NoError(t, errStat, "the directory must still exist")
			})
		}
	})

	t.Run("removes a sibling directory whose name starts like the current one", func(t *testing.T) {
		// "/path/to/go-screenplay-copy" contains the cwd "/path/to/go-screenplay"
		// as a string prefix but not as a directory: it must not be protected.
		sibling := currentFolder + "-sibling"
		require.NoError(t, os.Mkdir(sibling, fs.ModePerm))
		require.NoError(t, useTheFileSystem.ToRemoveDirectory(sibling))

		_, err := os.Stat(sibling)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("should be able to forget", func(t *testing.T) {
		require.NoError(t, useTheFileSystem.Forget())
	})

	t.Run("refuses to read a file larger than the buffered size", func(t *testing.T) {
		// A sparse file: the size is set without writing the bytes.
		filename := filepath.Join(t.TempDir(), "huge.bin")
		require.NoError(t, os.WriteFile(filename, nil, 0o600))
		require.NoError(t, os.Truncate(filename, screenplay.DefaultMaxBufferedSize+1))

		_, err := useTheFileSystem.ToRead(filename)
		require.EqualError(t, err, "the content of the file exceeds 64 MiB")

		require.NoError(t, os.Truncate(filename, screenplay.DefaultMaxBufferedSize))
		content, err := useTheFileSystem.ToRead(filename)
		require.NoError(t, err)
		assert.Len(t, content, screenplay.DefaultMaxBufferedSize)
	})

	t.Run("fails every operation on a nil receiver", func(t *testing.T) {
		var missing *ability.UseTheFileSystemAbility
		_, err := missing.ToRetrieveCurrentPath()
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		require.ErrorIs(t, missing.ToChangeDirectory("."), screenplay.ErrActorMissingAbility)
		_, err = missing.ToTellIfTheFileExists("notes.txt")
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		_, err = missing.ToTellIfTheDirectoryExists("notes")
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		_, err = missing.ToCreateTemporaryFile("notes-*.txt")
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		_, err = missing.ToCreateTemporaryFileWithContent("notes-*.txt", strings.NewReader("hello"))
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		require.ErrorIs(t, missing.ToCreateFile("notes.txt"), screenplay.ErrActorMissingAbility)
		require.ErrorIs(
			t,
			missing.ToCreateFileWithContent("notes.txt", strings.NewReader("hello")),
			screenplay.ErrActorMissingAbility,
		)
		require.ErrorIs(
			t,
			missing.ToOverwriteFileWithContent("notes.txt", strings.NewReader("hello")),
			screenplay.ErrActorMissingAbility,
		)
		require.ErrorIs(
			t,
			missing.ToAppendTo("notes.txt", strings.NewReader("hello")),
			screenplay.ErrActorMissingAbility,
		)
		_, err = missing.ToRead("notes.txt")
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		require.ErrorIs(t, missing.ToCreateDirectory("notes", 0o750), screenplay.ErrActorMissingAbility)
		_, err = missing.ToCreateTemporaryDirectory("notes-*")
		require.ErrorIs(t, err, screenplay.ErrActorMissingAbility)
		require.ErrorIs(t, missing.ToRemoveFile("notes.txt"), screenplay.ErrActorMissingAbility)
		require.ErrorIs(t, missing.ToRemoveDirectory("notes"), screenplay.ErrActorMissingAbility)
		assert.ErrorIs(t, missing.Forget(), screenplay.ErrActorMissingAbility)
	})
}

type failingReader struct{}

// Read always fails.
func (r *failingReader) Read(_ []byte) (int, error) {
	return 0, os.ErrInvalid
}
