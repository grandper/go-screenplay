package action

import "errors"

var (
	// ErrFilenameNotProvided is returned when no filename is provided.
	ErrFilenameNotProvided = errors.New("filename must be provided")

	// ErrDirectoryNameNotProvided is returned when no directory name is provided.
	ErrDirectoryNameNotProvided = errors.New("directory name must be provided")

	// ErrReaderNotProvided is returned when no reader is provided.
	ErrReaderNotProvided = errors.New("reader must be provided")

	// ErrContentNotProvided is returned when no content is provided.
	ErrContentNotProvided = errors.New("content must be provided")

	// ErrFileDoesNotExist is returned when the file to remove does not exist.
	ErrFileDoesNotExist = errors.New("does not exist")

	// ErrDirectoryDoesNotExist is returned when the directory to remove does not exist.
	ErrDirectoryDoesNotExist = errors.New("does not exist")
)
