# go-screenplay filesystem extension

This extension lets an actor drive the file system. It bundles the ability to use the
file system together with the actions, questions, and resolutions that go with it, so
you test file-system interactions with the same screenplay vocabulary you use
everywhere else.

## Packages
The symbols shown below live in the extension's sub-packages:

| Package | Import path | Provides |
|---------|-------------|----------|
| `ability` | `github.com/grandper/go-screenplay/extensions/filesystem/ability` | `UseTheFileSystem` |
| `action` | `github.com/grandper/go-screenplay/extensions/filesystem/action` | `Create`, `Remove`, `AppendTheText`/`AppendTheBytes`/`AppendTheContent`, `OverwriteTo`, `ChangeDirectoryTo` |
| `question` | `github.com/grandper/go-screenplay/extensions/filesystem/question` | `FileNamed`, `DirectoryNamed`, `ContentOfTheFileNamed` |
| `resolution` | `github.com/grandper/go-screenplay/extensions/filesystem/resolution` | `Exists` |

Assertions are written with the core screenplay packages (`action/see`,
`resolution/equal`, ...) alongside this extension's `resolution.Exists`. For readability
the examples below drop the package qualifiers on the extension constructors; import the
packages above and qualify them as `action.Create()`, `ability.UseTheFileSystem()`, and
so on (note that the extension's `action`/`question`/`resolution` packages share their
name with the core ones, so you may need an import alias when you use both in the same
file).

## The ability
This extension introduces a single ability: using the file system. Give it to an actor
with `WhoCan`:
```go
anActor := screenplay.ActorNamed("Boby").WhoCan(UseTheFileSystem())
```

## Actions

### Create files and directories
Create a file (the action fails if the file already exists):
```go
err := anActor.AttemptsTo(Create().TheFile("foobar.txt"))
```
Provide the content at creation time from an `io.Reader`, a byte slice, or a string:
```go
err := anActor.AttemptsTo(Create().TheFile("foobar.txt").Containing(strings.NewReader("Hello World")))
err := anActor.AttemptsTo(Create().TheFile("foobar.txt").ContainingBytes([]byte("Hello World")))
err := anActor.AttemptsTo(Create().TheFile("foobar.txt").ContainingTheText("Hello World"))
```
Create a directory (the action fails if the directory already exists):
```go
err := anActor.AttemptsTo(Create().TheDirectory("MyDocuments"))
```

### Create temporary files and directories
Create a temporary file from a pattern. Because its name is generated (a `*` in the
pattern is replaced with a random string), capture the final name with `AndSaveNameTo`:
```go
var filename string
err := anActor.AttemptsTo(Create().TheTemporaryFile("foobar-*.txt").AndSaveNameTo(&filename))
```
As with a regular file, you can provide the content from an `io.Reader`, a byte slice,
or a string:
```go
var filename string
err := anActor.AttemptsTo(Create().TheTemporaryFile("foobar-*.txt").Containing(strings.NewReader("Hello World")).AndSaveNameTo(&filename))
err := anActor.AttemptsTo(Create().TheTemporaryFile("foobar-*.txt").ContainingBytes([]byte("Hello World")).AndSaveNameTo(&filename))
err := anActor.AttemptsTo(Create().TheTemporaryFile("foobar-*.txt").ContainingTheText("Hello World").AndSaveNameTo(&filename))
```
Create a temporary directory the same way:
```go
var dirname string
err := anActor.AttemptsTo(Create().TheTemporaryDirectory("MyDocuments-*").AndSaveNameTo(&dirname))
```

### Append to a file
Append data to a file from an `io.Reader`, a byte slice, or a string:
```go
err := anActor.AttemptsTo(AppendTheContent(strings.NewReader("Hello World")).To("foobar.txt"))
err := anActor.AttemptsTo(AppendTheBytes([]byte("Hello World")).To("foobar.txt"))
err := anActor.AttemptsTo(AppendTheText("Hello World").To("foobar.txt"))
```

### Overwrite a file
Replace the whole content of a file from an `io.Reader`, a byte slice, or a string:
```go
err := anActor.AttemptsTo(OverwriteTo("foobar.txt").WithTheContent(strings.NewReader("Hello World")))
err := anActor.AttemptsTo(OverwriteTo("foobar.txt").WithTheBytes([]byte("Hello World")))
err := anActor.AttemptsTo(OverwriteTo("foobar.txt").WithTheText("Hello World"))
```

### Remove a file or a directory
```go
err := anActor.AttemptsTo(Remove().TheFile("foobar.txt"))
err := anActor.AttemptsTo(Remove().TheDirectory("MyDocuments"))
```

### Change the current directory
```go
err := anActor.AttemptsTo(ChangeDirectoryTo("MyDocuments"))
```

## Questions and resolutions

### Whether a file or directory exists
`FileNamed` and `DirectoryNamed` answer with a value you assert against with this
extension's `Exists` resolution:
```go
err := anActor.Should(see.The(FileNamed("foobar.txt")).Is(Exists()))
err := anActor.Should(see.The(DirectoryNamed("foobar")).Is(Exists()))
```
Because `see.The` also exposes negated verbs, you can assert the opposite with the same
resolution:
```go
err := anActor.Should(see.The(FileNamed("foobar.txt")).IsNot(Exists()))
```

### The content of a file
`ContentOfTheFileNamed` answers with the file's content as a string:
```go
err := anActor.Should(see.The(ContentOfTheFileNamed("foobar.txt")).Is(equal.To("Hello World")))
```
