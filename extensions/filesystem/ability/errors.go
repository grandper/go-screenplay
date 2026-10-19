package ability

import impl "github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/ability"

var (
	// ErrFileAlreadyExists is returned when the file already exists.
	ErrFileAlreadyExists = impl.ErrFileAlreadyExists
	// ErrSymbolicLink is returned when a write is refused because the path is a
	// symbolic link: the ability never writes through a symbolic link.
	ErrSymbolicLink = impl.ErrSymbolicLink
	// ErrProtectedDirectory is returned when a directory is not removed because
	// it is the root of the file system, the home directory, the current working
	// directory or a directory that contains it.
	ErrProtectedDirectory = impl.ErrProtectedDirectory
)
