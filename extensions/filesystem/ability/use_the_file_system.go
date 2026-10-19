package ability

import (
	impl "github.com/grandper/go-screenplay/v2/extensions/filesystem/internal/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// UseTheFileSystem represents the ability of an actor to use the file system.
//
// The ability has its own current directory, which starts as the working
// directory of the process. change.TheCurrentDirectoryTo moves this directory
// only, never the process's, so two actors, or two concurrent branches, never
// disturb each other. Every relative path given to the actions and the
// questions of the extension is resolved against it.
func UseTheFileSystem() screenplay.Ability {
	return impl.UseTheFileSystem()
}
