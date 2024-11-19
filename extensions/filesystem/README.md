# go-screenplay-filesystem

The filesystem extension provides new abilities, actions, questions, and resolutions to deal with the filesystem.

## New Abilities
This extension introduces a unique capability to interract with the file system.
To use the new capability use
```go
anActor := screeplay.ActorNamed("boby").WhoCan(UseTheFileSystem())
```

## New Actions

### Create files and directories
You create a file using
```go
err := anActor.AttemptsTo(Create().TheFile("foobar.txt")
```
Note that the creation will fail if the file already exists.

You can provide the content the while at the creation time using an `io.Reader`, bytes or a string.
```go
err := anActor.AttemptsTo(Create().TheFile("foobar.txt").Containing(strings.NewReader("Hello World")))
err := anActor.AttemptsTo(Create().TheFile("foobar.txt").ContainingBytes([]byte("Hello World"))
err := anActor.AttemptsTo(Create().TheFile("foobar.txt").ContainingTheText("Hello World"))
```

You can create a directory using
```go
err := anActor.AttemptsTo(Create().TheDirectory("MyDocuments"))
```
Note that the creation will fail if the directory already exists.

### Create temporary files and directories
You create a file using a pattern
```go
var filename string
err := anActor.AttemptsTo(Create().TheTemporaryFile("foobar-*.txt").AndSaveNameTo(&filename))
```
Note that since the filename is randomly generated, we need to store it using ".AndSaveNameTo".

You can provide the content the while at the creation time using an `io.Reader`, bytes or a string.
```go
var filename string
err := anActor.AttemptsTo(Create().TheTemporaryFile("foobar-*.txt").Containing(strings.NewReader("Hello World")).AndSaveNameTo(&filename))
err := anActor.AttemptsTo(Create().TheTemporaryFile("foobar-*.txt").ContainingBytes([]byte("Hello World").AndSaveNameTo(&filename))
err := anActor.AttemptsTo(Create().TheTemporaryFile("foobar-*.txt").ContainingTheText("Hello World").AndSaveNameTo(&filename))
```

You can create a temporary directory using
```go
err := anActor.AttemptsTo(Create().TheTemporaryDirectory("MyDocuments-*").AndSaveNameTo(&filename))
```
Note that since the directory name is randomly generated, we need to store it using ".AndSaveNameTo".

### Remove a file or a directory
You can remove a file or a directory using
```go
err := anActor.AttemptsTo(Remove().TheFile("foobar.txt"))
err := anActor.AttemptsTo(Remove().TheDirectory("MyDocuments"))
```

### Append to a file
You can append data to a file using an `io.Reader`, bytes or a string:
```go
err := anActor.AttemptsTo(AppendTheContent(strings.NewReader("Hello World")).To("foobar.txt"))
err := anActor.AttemptsTo(AppendTheBytes([]byte("Hello World")).To("foobar.txt"))
err := anActor.AttemptsTo(AppendTheText("Hello World").To("foobar.txt"))
```

### Overwrite a file
You can overwrite a file using an `io.Reader`, bytes or a string.
```go
err := anActor.AttemptsTo(OverwriteTo("foobar.txt").WithTheContent(strings.NewReader("Hello World")))
err := anActor.AttemptsTo(OverwriteTo("foobar.txt").WithTheBytes("Hello World"))
err := anActor.AttemptsTo(OverwriteTo("foobar.txt").WithTheText("Hello World"))
```

### Change the current directory
You can change the current directory using
```go
err := anActor.AttemptsTo(ChangeDirectoryTo("MyDocuments"))
```

## New Questions
You can check if a file or a directory exists using
```go
err := anActor.AttemptsTo(see.The(FileNamed("foobar.txt"), exists()))
err := anActor.AttemptsTo(see.The(DirectoryNamed("foobar"), exists()))
```
You can also assert what the content of a file is:
```go
err := anActor.AttemptsTo(see.The(ContentOfTheFileNamed("foobar.txt"), equals("Hello World")))
```
