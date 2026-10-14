# go-screenplay CLI extension

This extension lets an actor drive command-line programs. It bundles the ability to
run CLI commands together with the actions and questions that go with it, so you test
a command-line tool with the same screenplay vocabulary you use everywhere else.

## Packages
The symbols shown below live in the extension's sub-packages:

| Package | Import path | Provides |
|---------|-------------|----------|
| `ability` | `github.com/grandper/go-screenplay/v2/extensions/cli/ability` | `RunCLICommands`, `Command`, `Result` |
| `action` | `github.com/grandper/go-screenplay/v2/extensions/cli/action` | `RunTheCommand`, `Type` |
| `question` | `github.com/grandper/go-screenplay/v2/extensions/cli/question` | `EnvironmentVariableNamed`, `Responses` |
| `exitcode` | `github.com/grandper/go-screenplay/v2/extensions/cli/question/exitcode` | `Of` |
| `standard` | `github.com/grandper/go-screenplay/v2/extensions/cli/question/standard` | `OutputOf`, `ErrorOf` |

Assertions are written with the core screenplay packages (`action/see`,
`resolution/be`, `resolution/contain`, `resolution/length`, `question/last`, ...).
For readability the examples below drop the package qualifiers on the extension
constructors; import the packages above and qualify them as `action.RunTheCommand(...)`,
`ability.RunCLICommands()`, and so on. The extension's `action`/`question` packages share
their name with the core ones: when a file uses both, keep the core name and alias the
extension package with the `cli` prefix:
```go
import (
	"github.com/grandper/go-screenplay/v2/action"
	cliaction "github.com/grandper/go-screenplay/v2/extensions/cli/action"
)

err := anActor.AttemptsTo(action.Eventually(cliaction.RunTheCommand("ls")))
```

## The ability
This extension introduces a single ability: running CLI commands. Give it to an actor
with `WhoCan`:
```go
anActor := screenplay.ActorNamed("Boby").WhoCan(RunCLICommands())
```
The ability records the `Result` of every command it runs, so later questions can ask
about them.

An action or a question uses the ability in one chained statement, written in
`PerformAs` or `AnsweredBy`, never in a test:
```go
err := screenplay.UseAnAbilityOf(theActor).That(RunCLICommands()).ToRun(theActor.Context(), cmd)
```
The argument of `That` carries the type; `That` returns the actor's own instance,
or a nil pointer when the actor cannot run CLI commands, and the operation then
reports `screenplay.ErrActorMissingAbility`.

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
Pass environment variables, added to the environment the test inherits, one at
a time or as a map. A variable named like an inherited one overrides it:
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
err := anActor.Should(see.The(EnvironmentVariableNamed("PATH")).Should(be.EqualTo("/usr/local/bin:/usr/bin:/bin")))
```

### The responses
`Responses()` answers with every `*ability.Result` recorded so far, in the order the
commands were run. Because it answers with a slice, it composes with the core
`first`/`last`/`number` questions, and with the `exitcode` and `standard` questions
below to extract a specific field from a single response:
```go
err := anActor.Should(see.The(Responses()).ShouldHave(length.Exactly[[]*ability.Result](2)))
```

### The exit code
`exitcode.Of` wraps a question that answers with a `*Result` and answers with its
exit code (`0` means success):
```go
err := anActor.Should(see.The(exitcode.Of(last.Of(Responses()))).Should(be.EqualTo(0)))
```

### The standard output
`standard.OutputOf` answers with the standard output of a `*Result`:
```go
err := anActor.Should(see.The(standard.OutputOf(last.Of(Responses()))).Should(contain.TheBytes([]byte("Hello World"))))
```

### The standard error
`standard.ErrorOf` answers with the standard error of a `*Result`:
```go
err := anActor.Should(see.The(standard.ErrorOf(last.Of(Responses()))).Should(contain.TheBytes([]byte("unknown parameter '-x'"))))
```

## Reading a result directly
A `*ability.Result` exposes its parts through getters — `ExitCode()` (an `int`),
`StdOut()` and `StdErr()` (each a `[]byte`), and `Err()` (the error returned while
waiting for the command, if any):
```go
result, _ := anActor.AsksFor(last.Of(Responses())) // result is a *ability.Result

exitCode := result.ExitCode()
output := string(result.StdOut())
```
