# go-screenplay

`go-screenplay` is a package that brings the *screenplay* pattern to Golang.

First, let's give credits, where credits is due. This package takes some inspiration from several other libraries written in other languages:
- https://github.com/HamedStack/HamedStack.Playwright.Screenplay
- https://github.com/testla-project/testla-screenplay-playwright-js
- https://github.com/ScreenPyHQ/screenpy
- https://serenity-js.org/handbook/design/screenplay-pattern/

**Main features:**

- Implementation of the *Screenplay* Pattern in Golang.
- Utilities to use *Given-When-Then* instruction in tests.
- Natural-language-enabling syntactic sugar

## Installation

```bash
go get github.com/grandper/go-screenplay
```

## What is the Screenplay Pattern?
The Screenplay Pattern helps
- developing tests that are easy to maintain and that can keep up with a growing code base
that is constantly evolving  to add more features is challenging.
- minimizing the maintenance code, which is useful for developers and testers who try to treat their code for tests as production code.

### The Old Way: The Page Object Pattern
The Page Object Pattern is a simple UI pattern that provide a simple way to test UI code.
The main concept is to create page objects that contains the interaction for a given page.
The problem with the Page Object Pattern is that it lacks scalability because it breaks the Single Responsibility Principle of SOLID.

The basic structure of a test code for page object looks like
- model
  - Domain model classes
- pages
  - Page objects
- Steps
  - Tasks and assertions performed by the user.

### The Screenplay Pattern
This pattern was previously known as "the Journey Pattern" and was originally promoted. It is designed to create
tests from the point of view of the user. It uses a composition of objects that follows the SOLID principles to creates
easy to understand and maintainable tests.

It provides an efficient alternative to the so-called Page Object Pattern.

The pattern organizes the code into several components:
1. Element Locator provide element location details on the UI.
2. Actions present interaction implementation details
3. Tasks provide an object for users.
4. Questions provide information about what the user can see.
5. Actors represents the entity that takes actions and asks questions.

When using the Screenplay Pattern, the test code is usually organized as follows:
- model
  - Domain model classes
- tasks
  - Business tasks
- action
  - UI interactions (sub-tasks)
- pages
  - Page Objects and Elements
- questions
  - Objects used to query the application

## Main Concepts
In this section we go into more details about the different pieces of the Screenplay Pattern.

### Actor
The actor represents a user of your application. It is the glue to acts on actions, questions, and Resolution.
An actor is typically named to give some context about his role. This is also very useful when multiple actors
are expected to interact with the application for a given scenario to differentiate them.

To create an actor use
```go
reviewer := screenplay.ActorNamed("reviewer")
```
An actor can gain additional actions using Abilities.

### Abilities
Abilities brings new feature to your actors.
For example, you can make them use libraries, tools, and resources.

To give your actor an ability you can simply use
```go
neo := screenplay.ActorNamed("Neo").WhoCan(UseTerminalCommands())
```
Ability are then used to perform more actions. For example if you have
an ability that contains database sessions, then you can create an
action `QueryDatabase` that will be able to access the session of an
actor. You can also use ability in Questions.

Finally, abilities all implement a `Forget() error` interface that is
used to clean up the of what is inside the ability (e.g., closing
sessions). This method is automatically called when you use
```go
err := neo.Exit()
```

### Actor specific session data
Sometimes we need to store information in a step or rask and reuse it in a subsequent one.
To do so, each actor can call the `remember` and `recall` methods to store and retrieve
data in a key-value store.
```go
actor.Remember("the user id", 1234)
userID := actor.Recall("the user id").(int)
```
You can also hand `Remember` a question directly. The actor answers it and stores
the answer instead of the question. If the question fails to be answered, `nil` is
stored for the key:
```go
actor.Remember("the user id", UserID())
userID := actor.Recall("the user id").(int)
```
An actor can forget the data it has stored using
```go
actor.Forget("the user id")
```

### Cast
A cast provisions actors with abilities when they enter the stage.
It decouples the creation of actors from the abilities they need, so you only
define the setup once and every actor gets it automatically.

The simplest cast creates actors with no predefined abilities:
```go
cast := screenplay.CastOfStandardActors()
```
When all actors in a scenario share the same abilities, you can use:
```go
cast := screenplay.CastWhereEveryoneCan(MakeHTTPRequests(), UseDatabase())
```
For more control over how each actor is prepared you can supply a function:
```go
cast := screenplay.CastFunc(func(actor *screenplay.Actor) {
	actor.WhoCan(MakeHTTPRequests())
	actor.Remember("role", "reviewer")
})
```

### Stage
The stage manages all actors in a test scenario and tracks which one is
currently in the spotlight. It works together with a cast so that actors
are automatically prepared when they first appear.

To create a stage, provide it with a cast:
```go
stage := screenplay.SetTheStage(screenplay.CastWhereEveryoneCan(MakeHTTPRequests()))
```
Use `TheActorCalled` to retrieve an actor by name. If the actor does not exist
yet, the stage creates one and lets the cast prepare it. Calling it a second
time with the same name returns the same actor instance (name matching is
case-insensitive). The actor is also placed in the spotlight.
```go
adam := stage.TheActorCalled("Adam")
sameAdam := stage.TheActorCalled("adam") // returns the same actor
```
You can retrieve the actor currently in the spotlight at any time:
```go
actor, err := stage.TheActorInTheSpotlight()
```
To check whether any actor is currently on stage before accessing the spotlight:
```go
if stage.AnActorIsOnStage() {
	actor, _ := stage.TheActorInTheSpotlight()
	// use actor
}
```
At the end of a scenario, draw the curtain to make every actor exit and clean
up their abilities:
```go
err := stage.DrawTheCurtain()
```

### Actions
An actor will interact with the test under code by performing actions.
Actions can be used to set up your test or execute the operation that
you want to test.

Any action must implement the `Performable` interface:
```go
type Performable interface {
	PerformAs(actor *Actor) error
	String() string
}
```
Usually actions are created using fluent interface pattern to make the
code easier to read. For example, we can have a package `create` that
contains a method `Folder` that will create an action builder that
can be used as follows:
```go
action := create.Folder().Named("my_folder")
```
Now if we want to perform that action we can simply do one of the followings:
```go
bob := screenplay.ActorNamed("Bob")

// All the following functions will give you the name result.
err := bob.AttemptsTo(create.Folder().Named("my_folder"))
err := bob.WasAbleTo(create.Folder().Named("my_folder"))
err := bob.Does(create.Folder().Named("my_folder"))
err := bob.Did(create.Folder().Named("my_folder"))
err := bob.Will(create.Folder().Named("my_folder"))
err := bob.TriesTo(create.Folder().Named("my_folder"))
err := bob.TriedTo(create.Folder().Named("my_folder"))
err := bob.Tries(create.Folder().Named("my_folder"))
err := bob.Tried(create.Folder().Named("my_folder"))
err := bob.Shall(create.Folder().Named("my_folder"))
err := bob.Should(create.Folder().Named("my_folder"))
```
The reason for having all these alias functions available is providing you
the liberty to choose the word that makes the code of your tests easier to read.

### Tasks
It happens very often that we need to group several actions together.
We usually do this to avoid repeating ourselves or to give this group
of actions an explicit name.

For example, we may want to create a `Login` task that will group the following actions:
- `NavigateTo.TheHomePage()`
- `Click.On(TheLoginButton)`
- `FillIn.TheField(Username).With(username)`
- `FillIn.TheField(Password).With(password)`

To do so, we can use
```go
login := screenplay.TaskWhere("login to the application", 
	NavigateTo.TheHomePage(), 
	Click.On(TheLoginButton), 
	FillIn.TheField(Username).With(username), 
	FillIn.TheField(Password).With(password))
```
You can also wrap the call in a function to add parameters:
```go
func LoginAs(username, password string) screenplay.Performable {
	return screenplay.TaskWhere("login to the application", 
		NavigateTo.TheHomePage(), 
		Click.On(TheLoginButton), 
		FillIn.TheField(Username).With(username), 
		FillIn.TheField(Password).With(password))
}
```

Tasks implements the `Performable` interface just like any action so that an actor
can use tasks the same way:
```
err := bob.AttemptsTo(LoginAs("bob", "1234"))
```

### Questions and Resolutions
Tests are usually composed of three stages: the setup, the execution of some code,
and finally an assertion.
The counterpart of assertions in the Screenplay Pattern are questions and resolutions.
When an actor asks a question, the answer to this question can be asserted using a resolution.

A question is simply a struct that implements the following interface
```go
type Question interface {
	AnsweredBy(actor *Actor) (any, error)
	String() string
}
```
By passing the actor to the `AnsweredBy` method, we are able to use the actor `abilities` to answer the question.
The answer can be of any type.

The resolution is a struct that implements the following interface
```go
type Matcher func(obj any) (bool, error)

type Resolution interface {
	Resolve() Matcher
	String() string
}
```
As we can see here the resolution acts as a factory to instantiate a `Matcher`.
`Matcher` receives any type of objects and return `true` if that object
fulfills the expectation(s).

As mentioned before questions and resolutions works together: the actor asks a question to
retrieve an answer (i.e., a value of a given type) which is then asserted using the resolution.

Note that the resolution provides both which assertion is performed and the expected value.
For example, if we have a question `TitleOfThePage` that returns the title of an HTML page,
and a resolution that checks that the value is equal to "Some great title", then we can use
the action `see.The` to perform our test.
```go
err := theActor.Should(see.The(Text.OfThe(TitleOfThePage), is.EqualTo("Some great title")))
```

## Usage
Now that we have seen the foundation of the Screenplay Pattern and that we know how it works,
you probably realize that we need a collection of `abilities`, `actions`, `questions`, and `resolutions`.
In this section we will review what the library as already implemented for you.

#### Pausing, Stopping, and Waiting for Answers to Questions
You may want the test execution to top until the user hit the `Enter` key (mainly for debugging reasons).
In that case you can use the `Stop` action.
```go
err := theActor.Will(Stop())
```
You can also make the execution stop until some resolution is valid:
```go
err := theActor.Will(Stop().UntilThe(HomePage, contains.TheText("Hello World!")))
err := theActor.Will(Stop().UntilThe(CashAccount.Balance(), IsEqualTo(2000)))
```

You may also want to ask the execution for a given number of milliseconds, or seconds:
```go
err := theActor.AttemptsTo(PauseFor(10).Seconds().Because("the connection needs time to be setup"))
err := theActor.AttemptsTo(PauseFor(500).Milliseconds().Because("of some obscur reason"))
```
However, if you need to wait until something happens, you may want to use `Stop` instead to save time.

Finally, if you expect an action to fail repeatedly until finally it succeeds, you should use `Eventually`.
```go
err := adam.Should(Eventually(see.The(Text.OfThe(PageTitle), ContainsTheText("Hello World!"))))
err := adam.AttemptsTo(Eventually(Click.OnThe(SaveButton)).TryingEvery(100).Milliseconds())
err := adam.WasAbleTo(Eventually(CancelTheOrder).TryingFor(5).Seconds().PollingEvery(500).Milliseconds())
```

#### Logging
If you need to log the answer to a question, you can use the `Log` action:
```go
err := theActor.AttemptsTo(Log(HowManyBirdsAreInTheSky()))
err := theActor.AttemptsTo(Log(Number.Of(ItemsInTheList)))
```

#### Remembering answers to questions
When you need to ask a question and reuse its answer later, the `Remember` action
stores the answer in the actor's memory under a given key. You can then retrieve it
with `Recall` (see [Actor specific session data](#actor-specific-session-data)):
```go
err := theActor.AttemptsTo(Remember(StatusCode()).As("statusCode"))
statusCode := theActor.Recall("statusCode").(int)
```
By default the action fails if the question returns a `nil` answer. Use `AllowingNil`
to store it anyway:
```go
err := theActor.AttemptsTo(Remember(OptionalData()).As("data").AllowingNil())
```
Like any other action, it can be combined with `Eventually` to retry until an answer
is available:
```go
err := theActor.AttemptsTo(Eventually(Remember(UserProfile()).As("profile")).TryingFor(10).Seconds())
```

#### Working with multiple actions
Sometimes your actor needs to do a list of actions. In that case you can simply
provide a list of actions to your actor:
```go
err := theActor.AttemptsTo(DoThis(), DoThat())
```

#### Doing actions concurrently
When your actor needs to perform several actions in parallel, use `Concurrently`
(or its alias `Simultaneously`):
```go
err := theActor.AttemptsTo(Concurrently(DoThis(), DoThat(), DoSomethingElse()))
```
By default the actor waits for every action to complete and returns all the errors
that occurred, joined together with `errors.Join`. Calling `WaitingForAll` is optional
and makes this default explicit:
```go
err := theActor.AttemptsTo(Concurrently(DoThis(), DoThat()).WaitingForAll())
```
Use `StoppingOnError` to cancel the remaining actions as soon as one of them fails.
This mode relies on an `errgroup` and its context cancellation:
```go
err := theActor.AttemptsTo(Concurrently(DoThis(), DoThat()).StoppingOnError())
```
Use `IgnoringErrors` to run every action and discard the errors that occurred:
```go
err := theActor.AttemptsTo(Concurrently(DoThis(), DoThat()).IgnoringErrors())
```
Finally, you can cap the number of actions running at the same time with `WithLimit`:
```go
err := theActor.AttemptsTo(Concurrently(DoThis(), DoThat(), DoSomethingElse()).WithLimit(2).WaitingForAll())
```

#### Trying an action or doing an alternate actions
Sometimes you want the actor to try to do an action, and in case of failure do another one.
You can achieve this using `Either`:
```go
theActor.Will(Either(DoAction()).Or(DoDifferentAction())
theActor.Will(Either(DoAction()).Otherwise(DoDifferentAction()))
```

#### Observing things
The simplest way to ask a question is tu use the action `see`.
```go
err := theActor.Should(see.The(Text.OfThe(PageTile), StartsWith("Hello")))
```

It is possible to check if the actor sees any or all of a list of question-resolution pair.
```go
err := theActor.Should(see.AnyOf(
			profilWidget, contains.TheText("Adam"),
			loginForm, contains.TheText("Login")))
err := theActor.Should(see.AllOf(
			profilWidget, contains.TheText("Adam"),
			loginForm, contains.TheText("Login")))
```

You can check that all, any, or none of a set of questions match a resolution:
```go
err := theActor.Should(see.ThatAllOfThe(profilWidget, loginForm)(contains.TheText("Adam")))
err := theActor.Should(see.ThatAnyOfThe(profilWidget, loginForm)(contains.TheText("Adam")))
err := theActor.Should(see.ThatNoneOfThe(profilWidget, loginForm)(contains.TheText("Adam")))
```

#### Checking texts
You can check if a text starts with, ends with, or contains a string:
```go
err := theActor.Should(see.The(Text.OfThe(PageTile), StartsWith("Hello")))
err := theActor.Should(see.The(Text.OfThe(PageTile), EndsWith("World!")))
err := theActor.should(see.The(Text.OfThe(PageTitle), contains.TheText("lo Wor")))
```

You can check if a slice of bytes contains another slice of bytes:
```go
err := theActor.Should(see.The(BytesOfThe(LastResponse), contains.TheBytes([]byte("lo Wor"))))
```

You can match a text exactly:
```go
err := theActor.Should(see.The(Text.OfThe(PageTile), ReadsExactly("Hello World!")))
```

You can use regex to match a text:
```go
err := theActor.Should(see.The(Text.OfThe(PageTitle), Matches(`^Hello \w+`)))
pattern := regexp.MustCompile(`^Hello \w+`)
err := theActor.Should(see.The(Text.OfThe(PageTitle), Matches(pattern))
```

#### Checking numbers
You can assess the value of numbers as follows:
```go
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList)), is.EqualTo(0)))
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList)), is.LessThan(1)))
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList)), is.LessThanOrEqualTo(1)))
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList)), is.GreaterThan(1)))
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList)), is.GreaterThanOrEqualTo(1)))
delta := 25
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList)), is.CloseTo(101, delta)))
err := theActor.Should(see.The(Number.Of(items.InThe(TodoList)), is.InRange(1, 5)))
```

#### Checking collections
There are a couple of collection types in Golan such as slices and maps.

There exists a couple of resolutions to check slices:
```go
err := theActor.Should(see.The(List.OfAll(items.InThe(TodoList)), is.Empty()))
err := theActor.Should(see.The(items.InThe(TodoList), contains.TheItem("Add tests for the Go package")))
err := theActor.Should(see.The(items.InThe(TodoList), has.Length(5)))
err := theActor.should(see.The(Text.OfAll(items.InThe(TodoList)), contains.TheItem("by the end of the year")))
```

Here are a couple of options to test maps:
```go
err := theActor.Should(see.The(HeadersOf(TheLastResponse), contains.TheValue("application/json")))
err := theActor.Should(see.The(HeadersOf(TheLastResponse), contains.TheKey("Content-Type")))
err := theActor.Should(see.The(HeadersOf(TheLastResponse), contains.TheEntry("Content-Type", "application/json")))
```

When a question answers with a slice, you may want to assert against a single element of it
rather than the whole collection. The `first` and `last` questions wrap another question and
answer with, respectively, the first and last element of the slice it returns:
```go
err := theActor.Should(see.The(first.Of(items.InThe(TodoList)), ReadsExactly("Add tests for the Go package")))
err := theActor.Should(see.The(last.Of(items.InThe(TodoList)), ReadsExactly("by the end of the year")))
```
Both fail with `ErrAnswerIsNotASlice` if the wrapped question does not answer with a slice,
and with `ErrAnswerIsEmpty` if it answers with an empty one.

When a question answers with a struct, you may want to assert against its `Body` field rather
than the whole value. The `body` question wraps another question and answers with the `Body`
field of the struct it returns, whatever its type:
```go
err := theActor.Should(see.The(body.Of(TheLastResponse), is.EqualTo("created")))
```
It fails with `ErrAnswerIsNotAStruct` if the wrapped question does not answer with a struct (or a
pointer to one), and with `ErrFieldNotFound` if the struct has no `Body` field.

When a question answers with a slice or a map, you may want to assert against the number of
elements it contains. The `number` question wraps another question and answers with the count
of elements in the slice or map it returns:
```go
err := theActor.Should(see.The(number.Of(items.InThe(TodoList)), is.EqualTo(2)))
```
It fails with `ErrAnswerIsNotACollection` if the wrapped question does not answer with a slice or
a map.

#### Logical Operations
We can create the negation of a resolution:
```go
err := theActor.Should(see.The(PageTitle, is.Not(Visible())))
```

#### Creating an Anonymous Resolution
It happens sometimes that we need a one-off resolution without creating a dedicated struct.
We usually do this for simple or ad-hoc assertions that are not worth extracting into a reusable type.

For example, we may want to assert that a value satisfies a custom condition inline:
```go
isPositive := screenplay.ResolutionToSeeThatTheObject("is positive", func() screenplay.Matcher {
	return func(obj any) (bool, error) {
		n, ok := obj.(int)
		if !ok {
			return false, fmt.Errorf("expected an int, got %T", obj)
		}
		return n > 0, nil
	}
})
```
You can also wrap the call in a function to add parameters:
```go
func IsEqualTo(expected any) screenplay.Resolution {
	return screenplay.ResolutionToSeeThatTheObject(fmt.Sprintf("is equal to %v", expected), func() screenplay.Matcher {
		return func(obj any) (bool, error) {
			return obj == expected, nil
		}
	})
}
```

Anonymous resolutions implement the `Resolution` interface just like any other resolution so that an actor
can use them the same way:
```go
err := theActor.Should(see.The(CashAccount.Balance(), IsEqualTo(2000)))
```

#### Working with contexts
Contexts are hold by actors to share information across different part of a sequence of calls.
```go
ctx := context.Background()
actor := screenplay.ActorNamed("Alice").WithContext(ctx)
```
The context can later be retrieved using
```go
ctx := actor.Context()
```

## Using Gherkin-like syntax to write tests
You are maybe familiar with the Gherkin syntax that allow you to write
tests using the _Given-When-Then_ structure. Each keyword as a specific role:
* `Given` indicates that what follows is used to setup the test.
* `When` precedes the actions under test.
* `Then` precedes the assertions of the test.
The language also include the keyword `And` to avoid repeating any of the above
when multiple actions need to happen.

Here's an example:
> Given Adam's account has a balance of $2
> When Adam deposits $100
> Then the account balance is $102

The library provide some functions to reproduce the Gherkin flow in your test:
```go
adam := screenplay.ActorNamed("Adam")
screenplay.Given(adam).WasAbleTo(see.The(AccountBalance, is.Equal(2)))
screenplay.When(adam).AttemptsTo(Deposit(100).Dollars())
screenplay.Then(adam).Should(see.The(AccountBalance), is.EqualTo(102))
```

### Creating your own Performable/Task/Action
You can use `action.FromFunc` or the task tool.
To make the code more fluent, one can use builder methods:
```go
adam := screenplay.ActorNamed("Adam")
err := adam.AttemptsTo(FillIn().TheRegistrationForm().With(adamsData))
```

When the parameters of a task come from questions an actor must ask first, use
`action.FromFuncAndQuestions`. It asks each question, converts the answers to
the parameters of the function, and performs the returned task. The first
argument is a `fmt.Sprintf` format string used to describe the action — the
answers to the questions are passed as its format arguments:
```go
err := actor.AttemptsTo(
    action.FromFuncAndQuestions(
        "connect to %s on port %d",
        func(actor *screenplay.Actor, url string, port int) screenplay.Performable {
            return connect.To(url).On(port)
        },
        UrlQuestion{}, PortQuestion{},
    ),
)
```
The first parameter of the function must be `*screenplay.Actor` and the remaining
parameters must match, in order and type, the answers returned by the questions.

### Organizing Your Files
The organization of your test code is important. An easy way to organize the code is to group the different concept together.

- tests
  - features
    - feature1.go
      - Test for scenario 1
      - Test for scenario 2
    - feature2.go
  - actions
  - tasks
    - task1.go
    - task2.go
  - abilities
  - questions
  - resolutions

The features folder contains files for the different features. In each file, one can creates tests with different scenarios/test cases.

The tasks are different groups of actions. This way you have a collection of tasks that can be used to develop the tests of your features.

Then ability, actions, questions, and resolutions can also have their own folders.

## Extensions
You can find several extensions in the folder `extensions`. Each of them is used to extend the capability of the library to a specific use case.
- `http`: support API testing using REST requests.
- `cli`: support testing CLI applications.
- `filesystem`: support testing file system interactions.

## License
Licensed under MIT License.
