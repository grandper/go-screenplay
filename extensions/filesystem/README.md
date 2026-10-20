# go-screenplay filesystem extension

This extension lets an actor drive the file system. It bundles the ability to use the
file system together with the actions, questions, and resolutions that go with it, so
you test file-system interactions with the same screenplay vocabulary you use
everywhere else.

## Packages
The symbols shown below live in the extension's sub-packages:

| Package | Import path | Provides |
|---------|-------------|----------|
| `ability` | `github.com/grandper/go-screenplay/v2/extensions/filesystem/ability` | `UseTheFileSystem`, `ErrFileAlreadyExists`, `ErrSymbolicLink`, `ErrProtectedDirectory` |
| `action` | `github.com/grandper/go-screenplay/v2/extensions/filesystem/action` | `ErrFilenameNotProvided`, `ErrDirectoryNameNotProvided`, `ErrReaderNotProvided`, `ErrContentNotProvided`, `ErrFileDoesNotExist`, `ErrDirectoryDoesNotExist` |
| `create` | `github.com/grandper/go-screenplay/v2/extensions/filesystem/action/create` | `TheFile`, `TheTemporaryFile`, `TheDirectory`, `TheTemporaryDirectory` |
| `remove` | `github.com/grandper/go-screenplay/v2/extensions/filesystem/action/remove` | `TheFile`, `TheDirectory` |
| `write` | `github.com/grandper/go-screenplay/v2/extensions/filesystem/action/write` | `TheText`, `TheBytes`, `TheContent` |
| `overwrite` | `github.com/grandper/go-screenplay/v2/extensions/filesystem/action/overwrite` | `TheFile` |
| `change` | `github.com/grandper/go-screenplay/v2/extensions/filesystem/action/change` | `TheCurrentDirectoryTo` |
| `question` | `github.com/grandper/go-screenplay/v2/extensions/filesystem/question` | `FileNamed`, `DirectoryNamed`, `ContentOfTheFileNamed`, `File`, `Directory`, `NewFile`, `NewDirectory` |
| `resolution` | `github.com/grandper/go-screenplay/v2/extensions/filesystem/resolution` | `Exist` |
| `fixture` | `github.com/grandper/go-screenplay/v2/extensions/filesystem/question/fixture` | `CreateFakeExistingFile`, `CreateFakeMissingFile`, `CreateFakeExistingDirectory`, `CreateFakeMissingDirectory` |

Assertions are written with the core screenplay packages (`action/see`,
`resolution/be`, ...) alongside this extension's `resolution.Exist`. The actions are
written with their package, which is their verb: `create.TheFile("notes.txt")` reads
"create the file notes.txt". For readability the examples below drop the package
qualifiers on the other extension constructors; import the packages above and qualify
them as `ability.UseTheFileSystem()`, and so on. The extension's `question`/`resolution`
packages share their name with the core ones: when a file uses both, keep the core name
and alias the extension package with the `filesystem` prefix:
```go
import (
	"github.com/grandper/go-screenplay/v2/action/keep"
	"github.com/grandper/go-screenplay/v2/extensions/filesystem/action/remove"
)

err := anActor.AttemptsTo(keep.TryingTo(remove.TheFile("notes.txt")))
```

## The ability
This extension introduces a single ability: using the file system. Give it to an actor
with `WhoCan`:
```go
anActor := screenplay.ActorNamed("Boby").WhoCan(UseTheFileSystem())
```

`UseTheFileSystem()` returns a `screenplay.Ability`: a test gives it to an actor, and
cannot call its operations. The ability lives in `extensions/filesystem/internal/ability`,
so only an action or a question of the extension uses it, in one chained statement
written in `PerformAs` or `AnsweredBy`:
```go
err := screenplay.UseAnAbilityOf(theActor).That(impl.UseTheFileSystem()).ToCreateFile("notes.txt")
```
The argument of `That` carries the type; `That` returns the actor's own instance,
or a nil pointer when the actor cannot use the file system, and the operation
then reports `screenplay.ErrActorMissingAbility`.

The ability has its own current directory, which starts as the working directory of
the process. Every relative path given to the actions and the questions is resolved
against it (see [Change the current directory](#change-the-current-directory)).

When the actor exits, the ability has nothing to clean up: the files and directories
the actor created are left in place. Remove them with a cleanup task:
```go
anActor.HasIndependentCleanupTasks(remove.TheFile("foobar.txt"))
```

## Actions

### Create files and directories
Create a file (the action fails if the file already exists):
```go
err := anActor.AttemptsTo(create.TheFile("foobar.txt"))
// Boby attempts to create the file "foobar.txt"
```
Creating a file that already exists fails with `ErrFileAlreadyExists`, and the existing
file is left untouched:
```go
err := anActor.AttemptsTo(create.TheFile("foobar.txt"))
// file already exists
```
An empty name is reported when the actor performs the action, with
`ErrFilenameNotProvided`:
```go
err := anActor.AttemptsTo(create.TheFile(""))
// filename must be provided
```
A file is created with mode 0600, readable and writable by its owner only, and
no write ever follows a symbolic link: creating, appending to or overwriting a
path that is a symbolic link fails, with `ErrSymbolicLink` when appending or
overwriting:
```go
err := anActor.AttemptsTo(overwrite.TheFile("latest.log").WithTheText("Hello World"))
// the path is a symbolic link
```
Provide the content at creation time from an `io.Reader`, a byte slice, or a string:
```go
err := anActor.AttemptsTo(create.TheFile("foobar.txt").Containing(strings.NewReader("Hello World")))
err := anActor.AttemptsTo(create.TheFile("foobar.txt").ContainingBytes([]byte("Hello World")))
err := anActor.AttemptsTo(create.TheFile("foobar.txt").ContainingTheText("Hello World"))
// Boby attempts to create the file "foobar.txt" containing the text "Hello World"
```
The narration names the content as its method reads: `Containing` narrates "containing
the content", `ContainingBytes` "containing the bytes":
```go
err := anActor.AttemptsTo(create.TheFile("foobar.txt").ContainingBytes([]byte("Hello World")))
// Boby attempts to create the file "foobar.txt" containing the bytes "Hello World"
err := anActor.AttemptsTo(create.TheFile("foobar.txt").Containing(strings.NewReader("Hello World")))
// Boby attempts to create the file "foobar.txt" containing the content "Hello World"
```
`Containing` buffers the reader once, so the action can be performed more than once.
A reader holding more than 64 MiB is reported when the actor performs the action:
```go
err := anActor.AttemptsTo(create.TheFile("backup.tar").Containing(hugeArchive))
// failed to read the content: the content exceeds 64 MiB
```
Create a directory (the action fails if the directory already exists):
```go
err := anActor.AttemptsTo(create.TheDirectory("MyDocuments"))
// Boby attempts to create the directory "MyDocuments"
```
A directory is created with mode 0755. An empty name is reported when the actor
performs the action, with `ErrDirectoryNameNotProvided`:
```go
err := anActor.AttemptsTo(create.TheDirectory(""))
// directory name must be provided
```

### Create temporary files and directories
Create a temporary file from a pattern. Because its name is generated (the last `*` in
the pattern is replaced with a random string, which is appended when the pattern has no
`*`), capture the final name with `AndSaveNameTo`:
```go
var filename string
err := anActor.AttemptsTo(create.TheTemporaryFile("foobar-*.txt").AndSaveNameTo(&filename))
// Boby attempts to create the temporary file "foobar-*.txt"
```
A temporary file or directory is created in the temporary directory of the system,
never in the actor's current directory, and the name saved is its full path.
As with a regular file, you can provide the content from an `io.Reader`, a byte slice,
or a string:
```go
var filename string
err := anActor.AttemptsTo(create.TheTemporaryFile("foobar-*.txt").Containing(strings.NewReader("Hello World")).AndSaveNameTo(&filename))
err := anActor.AttemptsTo(create.TheTemporaryFile("foobar-*.txt").ContainingBytes([]byte("Hello World")).AndSaveNameTo(&filename))
err := anActor.AttemptsTo(create.TheTemporaryFile("foobar-*.txt").ContainingTheText("Hello World").AndSaveNameTo(&filename))
```
Create a temporary directory the same way:
```go
var dirname string
err := anActor.AttemptsTo(create.TheTemporaryDirectory("MyDocuments-*").AndSaveNameTo(&dirname))
// Boby attempts to create the temporary directory "MyDocuments-*"
```

### Write at the end of a file
Write data at the end of a file, without altering its existing content, from an
`io.Reader`, a byte slice, or a string:
```go
err := anActor.AttemptsTo(write.TheContent(strings.NewReader("Hello World")).AtTheEndOf("foobar.txt"))
err := anActor.AttemptsTo(write.TheBytes([]byte("Hello World")).AtTheEndOf("foobar.txt"))
err := anActor.AttemptsTo(write.TheText("Hello World").AtTheEndOf("foobar.txt"))
// Boby attempts to write the text "Hello World" at the end of "foobar.txt"
```
`TheBytes` and `TheContent` are narrated as they read:
```go
err := anActor.AttemptsTo(write.TheBytes([]byte("Hello World")).AtTheEndOf("foobar.txt"))
// Boby attempts to write the bytes "Hello World" at the end of "foobar.txt"
err := anActor.AttemptsTo(write.TheContent(strings.NewReader("Hello World")).AtTheEndOf("foobar.txt"))
// Boby attempts to write the content "Hello World" at the end of "foobar.txt"
```
The file must already exist. A missing file name or a nil reader is reported when the
actor performs the action, with `ErrFilenameNotProvided` or `ErrReaderNotProvided`:
```go
err := anActor.AttemptsTo(write.TheText("Hello World"))
// filename must be provided
err := anActor.AttemptsTo(write.TheContent(nil).AtTheEndOf("foobar.txt"))
// reader must be provided
```

### Overwrite a file
Replace the whole content of a file from an `io.Reader`, a byte slice, or a string:
```go
err := anActor.AttemptsTo(overwrite.TheFile("foobar.txt").WithTheContent(strings.NewReader("Hello World")))
err := anActor.AttemptsTo(overwrite.TheFile("foobar.txt").WithTheBytes([]byte("Hello World")))
err := anActor.AttemptsTo(overwrite.TheFile("foobar.txt").WithTheText("Hello World"))
// Boby attempts to overwrite the file "foobar.txt" with the text "Hello World"
```
`WithTheBytes` and `WithTheContent` are narrated as they read:
```go
err := anActor.AttemptsTo(overwrite.TheFile("foobar.txt").WithTheBytes([]byte("Hello World")))
// Boby attempts to overwrite the file "foobar.txt" with the bytes "Hello World"
err := anActor.AttemptsTo(overwrite.TheFile("foobar.txt").WithTheContent(strings.NewReader("Hello World")))
// Boby attempts to overwrite the file "foobar.txt" with the content "Hello World"
```
The file must already exist. Overwriting without a content is reported when the actor
performs the action, with `ErrContentNotProvided`:
```go
err := anActor.AttemptsTo(overwrite.TheFile("foobar.txt"))
// content must be provided
```

### Remove a file or a directory
```go
err := anActor.AttemptsTo(remove.TheFile("foobar.txt"))
// Boby attempts to remove the file "foobar.txt"
err := anActor.AttemptsTo(remove.TheDirectory("MyDocuments"))
// Boby attempts to remove the directory "MyDocuments"
```
Removing a file or a directory that does not exist fails with an error wrapping
`ErrFileDoesNotExist` or `ErrDirectoryDoesNotExist`:
```go
err := anActor.AttemptsTo(remove.TheFile("foobar.txt"))
// file "foobar.txt" does not exist
err := anActor.AttemptsTo(remove.TheDirectory("MyDocuments"))
// directory "MyDocuments" does not exist
```
Removing a directory removes its content too. The action refuses to remove the
root of the file system, the home directory, the current working directory and
any directory that contains it. The error wraps `ErrProtectedDirectory`:
```go
err := anActor.AttemptsTo(remove.TheDirectory(".."))
// Boby attempts to remove the directory "..": refused to remove a protected directory: ".." contains the current working directory
err := anActor.AttemptsTo(remove.TheDirectory("."))
// refused to remove a protected directory: "." is the current working directory
err := anActor.AttemptsTo(remove.TheDirectory("/"))
// refused to remove a protected directory: "/" is the root of the file system
```
The working directory of the process, and any directory that contains it, are
protected as well, even after the actor changed its current directory.

### Change the current directory
```go
err := anActor.AttemptsTo(change.TheCurrentDirectoryTo("MyDocuments"))
// Boby attempts to change the current directory to "MyDocuments"
```
The change applies to this actor only. Later relative paths of this actor
resolve against it, and the working directory of the process is never changed:
```go
err := anActor.AttemptsTo(
    change.TheCurrentDirectoryTo("MyDocuments"),
    create.TheFile("notes.txt"),      // creates MyDocuments/notes.txt
)
```
Changing to a path that is not a directory fails, and an empty name is reported with
`ErrDirectoryNameNotProvided`:
```go
err := anActor.AttemptsTo(change.TheCurrentDirectoryTo("foobar.txt"))
// chdir foobar.txt: not a directory
err := anActor.AttemptsTo(change.TheCurrentDirectoryTo(""))
// directory name must be provided
```

## Questions and resolutions

### Whether a file or directory exists
`FileNamed` and `DirectoryNamed` answer with a `*question.File` and a
`*question.Directory` you assert against with this extension's `Exist` resolution.
Go cannot infer which of the two a resolution is for, so the type is given explicitly:
```go
err := anActor.Should(see.The(FileNamed("foobar.txt")).Should(Exist[*question.File]()))
err := anActor.Should(see.The(DirectoryNamed("foobar")).Should(Exist[*question.Directory]()))
```
Because `see.The` also exposes negated verbs, you can assert the opposite with the same
resolution:
```go
err := anActor.Should(see.The(FileNamed("foobar.txt")).ShouldNot(Exist[*question.File]()))
```
When the file is missing, the failure reads:
```go
err := anActor.Should(see.The(FileNamed("foobar.txt")).Should(Exist[*question.File]()))
// Boby failed to see that the file named "foobar.txt" should exist: it does not exist
```
`FileNamed` answers with a file whose path is a regular file, and `DirectoryNamed` with
a directory whose path is a directory: a directory named `foobar.txt` is not a file.

### Reading a file or a directory directly
A `*question.File` and a `*question.Directory` expose their parts through getters —
`Dir()` and `Name()` (the path split as it was written, not resolved) and `Exists()`:
```go
file, _ := anActor.AsksFor(FileNamed("reports/2024.csv")) // file is a *question.File

dir := file.Dir()       // "reports"
name := file.Name()     // "2024.csv"
exists := file.Exists()
```

### The content of a file
`ContentOfTheFileNamed` answers with the file's content as a `[]byte`:
```go
err := anActor.Should(see.The(ContentOfTheFileNamed("foobar.txt")).Should(be.DeepEqualTo([]byte("Hello World"))))
content, err := anActor.AsksFor(ContentOfTheFileNamed("foobar.txt")) // content is a []byte
```
A file larger than 64 MiB is refused, so a big fixture never exhausts the memory of the
test:
```go
content, err := anActor.AsksFor(ContentOfTheFileNamed("backup.tar"))
// the content of the file exceeds 64 MiB
```

## Fakes for your own questions and resolutions
The `fixture` package creates a `*question.File` or a `*question.Directory` without
touching the file system, to test code that receives the answers of `FileNamed` and
`DirectoryNamed`, such as a resolution of your own. `question.NewFile` and
`question.NewDirectory` do the same with the existence given as an argument:
```go
notes := fixture.CreateFakeExistingFile("notes.txt")
draft := fixture.CreateFakeMissingFile("draft.txt")
documents := fixture.CreateFakeExistingDirectory("MyDocuments")
build := fixture.CreateFakeMissingDirectory("build")

err := Exist[*question.File]().Inspect(draft)
// it does not exist
```
