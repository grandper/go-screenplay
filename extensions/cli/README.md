# go-screenplay-cli

The cli extension provides new abilities, actions, questions, and resolutions to deal with the command line.

## New Abilities
This extension introduces a unique capability to interact with the command line.
To use the new capability use
```go
anActor := screenplay.ActorNamed("Boby").WhoCan(RunCLICommands())
```

## New Actions

### Run a command
You can run a command using the following action:
```go
err := anActor.AttemptsTo(RunTheCommand("echo", "'Hello World'"))
```

Commands can be run in a specific working directory:
```go
err := anActor.AttemptsTo(RunTheCommand("ls").InTheWorkingDirectory("/home/boby"))
```

It is also possible to run a command interactively, which means that you can provide input to the command line:
```go
err := anActor.AttemptsTo(RunTheCommand("isprime").Interactively()
```

### Type input to the command line
You can type input to the command line using the following action:
```go
err := anActor.AttemptsTo(Type("42"))
```
You can also simulate pressing the Enter key:
```go
err := anActor.AttemptsTo(Type("42").AndPressEnter())
```
You can format the input using fmt.Sprintf style formatting:
```go
err := anActor.AttemptsTo(Type("%s %s", "Hello", "World"))
```

## New Questions

### Environment Variables
You can ask for the value of an environment variable using the following question:
```go
err := anActor.Should(see.The(EnvironmentVariableNamed("PATH"), is.EqualTo("/usr/local/bin:/usr/bin:/bin")))
```

### The Error Code
You can request information about the error code of the last response.
```go
err := anActor.Should(see.The(ErrorCodeOfTheLastResponse(), is.EqualTo(0)))
```

### The stdout
You can request information about the standard output of the last response.
```go
err := anActor.Should(see.The(StandardOutputOfTheLastResponse(), contains.TheText("Hello World")))
```

### The stderr
You can request information about the standard error of the last response.
```go
err := anActor.Should(see.The(StandardErrorOfTheLastResponse(), contains.TheText("unknown parameter '-x'")))
```
