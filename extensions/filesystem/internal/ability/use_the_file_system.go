package ability

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/verdict"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

var (
	// ErrFileAlreadyExists is returned when the file already exists.
	ErrFileAlreadyExists = errors.New("file already exists")
	// ErrSymbolicLink is returned when a write is refused because the path is a
	// symbolic link: the ability never writes through a symbolic link.
	ErrSymbolicLink = errors.New("the path is a symbolic link")
	// ErrProtectedDirectory is returned when ToRemoveDirectory refuses to remove
	// the root of the file system, the home directory, the current working
	// directory or a directory that contains it.
	ErrProtectedDirectory = errors.New("refused to remove a protected directory")
)

// createdFileMode is the permission of every file the ability creates: readable
// and writable by its owner only, like a temporary file from os.CreateTemp.
const createdFileMode = 0600

// UseTheFileSystem represents the ability of an actor to use the file system.
func UseTheFileSystem() *UseTheFileSystemAbility {
	return &UseTheFileSystemAbility{}
}

// UseTheFileSystemAbility represents the ability of an actor to use the file system.
//
// The ability has its own current directory, which starts as the working
// directory of the process. ToChangeDirectory moves this directory only, never
// the process's, so two actors, or two concurrent branches, never disturb each
// other. Every relative path given to the ability is resolved against it.
type UseTheFileSystemAbility struct {
	mutex sync.RWMutex
	// currentDirectory is the ability's own current directory; empty until
	// ToChangeDirectory is called, meaning the working directory of the process.
	currentDirectory string
}

// ToRetrieveCurrentPath returns the ability's current directory.
func (a *UseTheFileSystemAbility) ToRetrieveCurrentPath() (string, error) {
	if a == nil {
		return "", screenplay.ErrActorMissingAbility
	}
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	if a.currentDirectory != "" {
		return a.currentDirectory, nil
	}
	return os.Getwd()
}

// ToChangeDirectory changes the ability's current directory. The working directory
// of the process is left untouched.
func (a *UseTheFileSystemAbility) ToChangeDirectory(dir string) error {
	if a == nil {
		return screenplay.ErrActorMissingAbility
	}
	path, err := a.resolve(dir)
	if err != nil {
		return err
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("chdir %s: not a directory", dir)
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.currentDirectory = path
	return nil
}

// resolve returns the clean absolute path of name: an absolute name as is, a
// relative name under the ability's current directory.
func (a *UseTheFileSystemAbility) resolve(name string) (string, error) {
	if filepath.IsAbs(name) {
		return filepath.Clean(name), nil
	}
	current, err := a.ToRetrieveCurrentPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(current, name), nil
}

// ToTellIfTheFileExists returns true if the file exists.
func (a *UseTheFileSystemAbility) ToTellIfTheFileExists(filename string) (bool, error) {
	if a == nil {
		return false, screenplay.ErrActorMissingAbility
	}
	path, err := a.resolve(filename)
	if err != nil {
		return false, err
	}
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// ToTellIfTheDirectoryExists returns true if the directory exists.
func (a *UseTheFileSystemAbility) ToTellIfTheDirectoryExists(directory string) (bool, error) {
	if a == nil {
		return false, screenplay.ErrActorMissingAbility
	}
	path, err := a.resolve(directory)
	if err != nil {
		return false, err
	}
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// ToCreateTemporaryFile creates a new temporary file whose name is built from
// the pattern, and returns that name. A random string is added to the end of
// the pattern, or replaces its last "*". It is the caller's responsibility to
// remove the file when it is no longer needed.
func (a *UseTheFileSystemAbility) ToCreateTemporaryFile(pattern string) (string, error) {
	if a == nil {
		return "", screenplay.ErrActorMissingAbility
	}
	file, err := a.createTemporaryFile(pattern)
	if err != nil {
		return "", err
	}
	return file.Name(), file.Close()
}

// ToCreateTemporaryFileWithContent creates a new file with the specified name and copies the content from the provided reader.
func (a *UseTheFileSystemAbility) ToCreateTemporaryFileWithContent(pattern string, reader io.Reader) (string, error) {
	if a == nil {
		return "", screenplay.ErrActorMissingAbility
	}
	file, err := a.createTemporaryFile(pattern)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = io.Copy(file, reader)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

// createTemporaryFile creates the file in the temporary directory of the
// system. os.CreateTemp adds a random suffix to the pattern, so the name never
// collides with an existing file: no check against another directory is made.
func (a *UseTheFileSystemAbility) createTemporaryFile(pattern string) (*os.File, error) {
	return os.CreateTemp("", pattern)
}

// ToCreateFile creates a new file with the specified name.
func (a *UseTheFileSystemAbility) ToCreateFile(filename string) error {
	if a == nil {
		return screenplay.ErrActorMissingAbility
	}
	file, err := a.createFile(filename)
	if err != nil {
		return err
	}
	return file.Close()
}

// ToCreateFileWithContent creates a new file with the specified name and copies the content from the provided reader.
func (a *UseTheFileSystemAbility) ToCreateFileWithContent(filename string, reader io.Reader) error {
	if a == nil {
		return screenplay.ErrActorMissingAbility
	}
	file, err := a.createFile(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}
	return nil
}

// createFile creates the file in one atomic call: O_EXCL makes the kernel refuse
// a path that already exists, so no file appearing between a check and the
// creation can be truncated, and a symbolic link on the final component,
// dangling or not, is refused as well.
func (a *UseTheFileSystemAbility) createFile(filename string) (*os.File, error) {
	path, err := a.resolve(filename)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, createdFileMode)
	if errors.Is(err, os.ErrExist) {
		return nil, ErrFileAlreadyExists
	}
	if err != nil {
		return nil, err
	}
	return file, nil
}

// refuseSymbolicLink fails with ErrSymbolicLink when the path is a symbolic
// link, so that a write never lands somewhere else than the named path.
func refuseSymbolicLink(path string) error {
	info, err := os.Lstat(path)
	if err == nil && info.Mode()&os.ModeSymlink != 0 {
		return ErrSymbolicLink
	}
	return nil
}

// ToOverwriteFileWithContent overwrites the content of the specified file with the content from the provided reader.
func (a *UseTheFileSystemAbility) ToOverwriteFileWithContent(filename string, reader io.Reader) error {
	if a == nil {
		return screenplay.ErrActorMissingAbility
	}
	path, err := a.resolve(filename)
	if err != nil {
		return err
	}
	if err = refuseSymbolicLink(path); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, createdFileMode)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}
	return nil
}

// ToAppendTo appends content from the provided reader to the specified file.
func (a *UseTheFileSystemAbility) ToAppendTo(filename string, reader io.Reader) error {
	if a == nil {
		return screenplay.ErrActorMissingAbility
	}
	path, err := a.resolve(filename)
	if err != nil {
		return err
	}
	if err = refuseSymbolicLink(path); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, createdFileMode)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}
	return nil
}

// ToRead reads the content of the specified file and returns it as a byte slice.
// A file larger than screenplay.DefaultMaxBufferedSize is refused, so a big
// fixture never exhausts the memory of the test process.
func (a *UseTheFileSystemAbility) ToRead(filename string) ([]byte, error) {
	if a == nil {
		return nil, screenplay.ErrActorMissingAbility
	}
	path, err := a.resolve(filename)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(io.LimitReader(file, screenplay.DefaultMaxBufferedSize+1))
	if err != nil {
		return nil, err
	}
	if len(content) > screenplay.DefaultMaxBufferedSize {
		return nil, fmt.Errorf(
			"the content of the file exceeds %d MiB",
			screenplay.DefaultMaxBufferedSize>>20, //nolint:mnd // shifting by 20 converts bytes to MiB
		)
	}
	return content, nil
}

// ToCreateDirectory creates a new directory with the specified name and permission bits.
func (a *UseTheFileSystemAbility) ToCreateDirectory(name string, perm os.FileMode) error {
	if a == nil {
		return screenplay.ErrActorMissingAbility
	}
	path, err := a.resolve(name)
	if err != nil {
		return err
	}
	return os.Mkdir(path, perm)
}

// ToCreateTemporaryDirectory creates a new temporary directory whose name is
// built from the pattern, and returns that name. A random string is added to
// the end of the pattern, or replaces its last "*".
func (a *UseTheFileSystemAbility) ToCreateTemporaryDirectory(pattern string) (string, error) {
	if a == nil {
		return "", screenplay.ErrActorMissingAbility
	}
	return os.MkdirTemp("", pattern)
}

// ToRemoveFile removes the named file or (empty) directory.
func (a *UseTheFileSystemAbility) ToRemoveFile(filename string) error {
	if a == nil {
		return screenplay.ErrActorMissingAbility
	}
	path, err := a.resolve(filename)
	if err != nil {
		return err
	}
	return os.Remove(path)
}

// ToRemoveDirectory removes a directory and its content. It refuses, with
// ErrProtectedDirectory, to remove the root of the file system, the home
// directory, the current working directory or a directory that contains it.
func (a *UseTheFileSystemAbility) ToRemoveDirectory(name string) error {
	if a == nil {
		return screenplay.ErrActorMissingAbility
	}
	path, err := a.resolve(name)
	if err != nil {
		return err
	}
	if err = a.refuseProtectedDirectory(name, path); err != nil {
		return err
	}
	return os.RemoveAll(path)
}

// refuseProtectedDirectory fails with ErrProtectedDirectory when removing the
// resolved directory would be catastrophic: the root of the file system, the
// home directory, the ability's current directory or the working directory of
// the process, or a directory that contains one of those two.
func (a *UseTheFileSystemAbility) refuseProtectedDirectory(name, dir string) error {
	if dir == filepath.Dir(dir) {
		return fmt.Errorf("%w: %s is the root of the file system", ErrProtectedDirectory, verdict.Describe(name))
	}

	if home, err := os.UserHomeDir(); err == nil && dir == filepath.Clean(home) {
		return fmt.Errorf("%w: %s is the home directory", ErrProtectedDirectory, verdict.Describe(name))
	}

	current, err := a.ToRetrieveCurrentPath()
	if err != nil {
		return err
	}
	if dir == current {
		return fmt.Errorf("%w: %s is the current working directory", ErrProtectedDirectory, verdict.Describe(name))
	}
	if strings.HasPrefix(current, dir+string(filepath.Separator)) {
		return fmt.Errorf(
			"%w: %s contains the current working directory",
			ErrProtectedDirectory,
			verdict.Describe(name),
		)
	}

	processDirectory, err := os.Getwd()
	if err != nil {
		return err
	}
	if dir == processDirectory {
		return fmt.Errorf(
			"%w: %s is the working directory of the process",
			ErrProtectedDirectory,
			verdict.Describe(name),
		)
	}
	if strings.HasPrefix(processDirectory, dir+string(filepath.Separator)) {
		return fmt.Errorf(
			"%w: %s contains the working directory of the process",
			ErrProtectedDirectory,
			verdict.Describe(name),
		)
	}

	return nil
}

// Forget has nothing to clean up: the files and directories the actor created
// are left in place. On a nil receiver, the actor does not have the ability, it
// fails with screenplay.ErrActorMissingAbility.
func (a *UseTheFileSystemAbility) Forget() error {
	if a == nil {
		return screenplay.ErrActorMissingAbility
	}
	return nil
}

// UseTheFileSystemAbility implements the screenplay.Ability interface.
var _ screenplay.Ability = (*UseTheFileSystemAbility)(nil)
