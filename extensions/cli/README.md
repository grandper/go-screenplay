# go-screenplay CLI extension

This extension lets an actor drive command-line programs. It bundles the ability to
run CLI commands together with the actions and questions that go with it, so you test
a command-line tool with the same screenplay vocabulary you use everywhere else.

## Packages
The symbols shown below live in the extension's sub-packages:

| Package | Import path | Provides |
|---------|-------------|----------|
| `ability` | `github.com/grandper/go-screenplay/extensions/cli/ability` | `RunCLICommands`, `Command`, `Result` |
| `action` | `github.com/grandper/go-screenplay/extensions/cli/action` | `RunTheCommand`, `Type` |
| `question` | `github.com/grandper/go-screenplay/extensions/cli/question` | `EnvironmentVariableNamed`, `Responses` |
| `errorcode` | `github.com/grandper/go-screenplay/extensions/cli/question/errorcode` | `Of` |
| `standard` | `github.com/grandper/go-screenplay/extensions/cli/question/standard` | `OutputOf`, `ErrorOf` |

Assertions are written with the core screenplay packages (`action/see`,
`resolution/equal`, `resolution/contain`, `resolution/length`, `question/last`, ...).
For readability the examples below drop the package qualifiers on the extension
constructors; import the packages above and qualify them as `action.RunTheCommand(...)`,
`ability.RunCLICommands()`, and so on (note that the extension's `action`/`question`
packages share their name with the core ones, so you may need an import alias when you
use both in the same file).

## The ability
This extension introduces a single ability: running CLI commands. Give it to an actor
with `WhoCan`:
```go
anActor := screenplay.ActorNamed("Boby").WhoCan(RunCLICommands())
```
The ability records the `Result` of every command it runs, so later questions can ask
about them.

## Actions

### Running a command
Run a command with `RunTheCommand`, passing the program name and its arguments:
```go
err := anActor.AttemptsTo(RunTheCommand("echo", "Hello World"))
```
Run it in a specific working directory:
```go
err := anActor.AttemptsTo(RunTheCommand("ls").InTheWorkingDirectory("/home/boby"))
```
Pass environment variables, one at a time or as a map:
```go
err := anActor.AttemptsTo(RunTheCommand("printenv", "TOKEN").WithEnvVar("TOKEN", "s3cret"))
err := anActor.AttemptsTo(RunTheCommand("printenv").WithEnv(map[string]string{
    "TOKEN": "s3cret",
    "DEBUG": "1",
}))
```
Run a command interactively so you can feed it input with the `Type` action:
```go
err := anActor.AttemptsTo(RunTheCommand("isprime").Interactively())
```

### Typing input to a command
Type input to a command started with `Interactively`. `Type` formats its argument
with `fmt.Sprintf`-style formatting:
```go
err := anActor.AttemptsTo(Type("42"))
err := anActor.AttemptsTo(Type("%s %s", "Hello", "World"))
```
Simulate pressing the Enter key after the input with `AndPressEnter`:
```go
err := anActor.AttemptsTo(Type("42").AndPressEnter())
```

## Questions

### Environment variables
Ask for the value of an environment variable:
```go
err := anActor.Should(see.The(EnvironmentVariableNamed("PATH")).Is(equal.To("/usr/local/bin:/usr/bin:/bin")))
```

### The responses
`Responses()` answers with every `*ability.Result` recorded so far, in the order the
commands were run. Because it answers with a slice, it composes with the core
`first`/`last`/`number` questions, and with the `errorcode` and `standard` questions
below to extract a specific field from a single response:
```go
err := anActor.Should(see.The(Responses()).Has(length.Of(2)))
```

### The exit code
`errorcode.Of` wraps a question that answers with a `*Result` and answers with its
exit code (`0` means success):
```go
err := anActor.Should(see.The(errorcode.Of(last.Of(Responses()))).Is(equal.To(0)))
```

### The standard output
`standard.OutputOf` answers with the standard output of a `*Result`:
```go
err := anActor.Should(see.The(standard.OutputOf(last.Of(Responses()))).Does(contain.TheText("Hello World")))
```

### The standard error
`standard.ErrorOf` answers with the standard error of a `*Result`:
```go
err := anActor.Should(see.The(standard.ErrorOf(last.Of(Responses()))).Does(contain.TheText("unknown parameter '-x'")))
```

## Reading a result directly
A `*ability.Result` exposes its parts through getters — `ExitCode()` (an `int`),
`StdOut()` and `StdErr()` (each a `[]byte`), and `Err()` (the error returned while
waiting for the command, if any):
```go
answer, _ := last.Of(Responses()).AnsweredBy(anActor)
result := answer.(*ability.Result)

exitCode := result.ExitCode()
output := string(result.StdOut())
```
