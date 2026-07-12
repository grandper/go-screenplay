# go-screenplay HTTP extension

This extension lets an actor drive REST APIs. It bundles the ability to make HTTP
requests together with the actions, questions, and resolutions that go with it, so
you test an API with the same screenplay vocabulary you use everywhere else.

## Packages
The symbols shown below live in the extension's sub-packages:

| Package | Import path | Provides |
|---------|-------------|----------|
| `ability` | `github.com/grandper/go-screenplay/extensions/http/ability` | `MakeHTTPRequests`, `HTTPResponse`, `Credential` |
| `action` | `github.com/grandper/go-screenplay/extensions/http/action` | `SendHTTPRequest` (and the per-method shortcuts), `AddHeader`, `AddHeaders`, `SetHeader`, `SetHeaders` |
| `question` | `github.com/grandper/go-screenplay/extensions/http/question` | `Responses` |
| `status` | `github.com/grandper/go-screenplay/extensions/http/question/status` | `CodeOf` |

Assertions are written with the core screenplay packages (`action/see`,
`resolution/equal`, `resolution/contain`, `question/last`, ...). For readability the
examples below drop the package qualifiers on the extension constructors; import the
packages above and qualify them as `action.SendGetRequest()`, `ability.MakeHTTPRequests()`,
and so on (note that the extension's `action`/`question` packages share their name with
the core ones, so you may need an import alias when you use both in the same file).

## The ability
This extension introduces a single ability: making HTTP requests. Give it to an actor
with `WhoCan`:
```go
anActor := screenplay.ActorNamed("Wanda").WhoCan(MakeHTTPRequests())
```
The ability keeps the headers set for the session and records every response it
receives, so later questions can ask about them.

## Actions

### Headers
Add a header (keeping any header already set), or set the headers (replacing the
whole set) for every subsequent request of the session:
```go
err := anActor.AttemptsTo(AddHeader("Content-Type", "application/json"))
err := anActor.AttemptsTo(AddHeaders(
    "Content-Type", "application/json",
    "Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686"))
err := anActor.AttemptsTo(SetHeader("Content-Type", "application/json"))
err := anActor.AttemptsTo(SetHeaders(
    "Content-Type", "application/json",
    "Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686"))
```
`AddHeaders` and `SetHeaders` take an even number of arguments, read as
`key, value, key, value, ...`.

If a header holds a secret you do not want in the narration, mark the action secret
with `Secretly` (or its longer alias `WhichShouldBeKeptSecret`):
```go
err := anActor.AttemptsTo(AddHeader("Authorization", "Bearer "+token).Secretly())
err := anActor.AttemptsTo(SetHeader("Authorization", "Bearer "+token).WhichShouldBeKeptSecret())
```

### Sending requests
Send a request with `SendHTTPRequest`, choosing the method, the URL (`To`), and
optionally a body (`WithBody`):
```go
err := anActor.AttemptsTo(SendHTTPRequest(http.MethodGet).To("http://www.example.com"))
err := anActor.AttemptsTo(SendHTTPRequest(http.MethodPost).To("http://www.example.com").WithBody(body))
```
For readability there is a shortcut per method:
```go
err := anActor.AttemptsTo(SendDeleteRequest().To("http://www.example.com"))
err := anActor.AttemptsTo(SendGetRequest().To("http://www.example.com"))
err := anActor.AttemptsTo(SendHeadRequest().To("http://www.example.com"))
err := anActor.AttemptsTo(SendOptionsRequest().To("http://www.example.com"))
err := anActor.AttemptsTo(SendPatchRequest().To("http://www.example.com").WithBody(body))
err := anActor.AttemptsTo(SendPostRequest().To("http://www.example.com").WithBody(body))
err := anActor.AttemptsTo(SendPutRequest().To("http://www.example.com").WithBody(body))
```
`WithBody` takes an `io.Reader`. It buffers the reader so the request can be
described in the narration and sent (and re-sent) without the body being consumed.

Use basic authentication with `WithCredential` (or its alias `WithAuth`):
```go
err := anActor.AttemptsTo(SendGetRequest().To("http://www.example.com").WithCredential("username", "password"))
err := anActor.AttemptsTo(SendGetRequest().To("http://www.example.com").WithAuth("username", "password"))
```
A request can be sent secretly, so neither its body nor headers appear in the
narration:
```go
err := anActor.AttemptsTo(SendPostRequest().To("http://www.example.com").WithBody(body).Secretly())
err := anActor.AttemptsTo(SendPostRequest().To("http://www.example.com").WithBody(body).WhichShouldBeKeptSecret())
```

## Questions

### The responses
After sending one or more requests, `Responses()` answers with every
`*ability.HTTPResponse` the actor received, in the order the requests were sent.
Because it answers with a slice, it composes with the core `first`/`last`/`number`
questions to focus on a single response:
```go
err := anActor.Should(see.The(number.Of(Responses())).Is(equal.To(2)))
```

### The status code
`status.CodeOf` wraps a question that answers with an `*HTTPResponse` and answers
with its status code, so you assert on it like any number:
```go
err := anActor.Should(see.The(status.CodeOf(last.Of(Responses()))).Is(equal.To(200)))
```

### Reading a response directly
An `*ability.HTTPResponse` exposes its parts through getters — `Body()` (a
`string`), `Headers()` (a `map[string]string`), and `StatusCode()` (an `int`):
```go
answer, _ := last.Of(Responses()).AnsweredBy(anActor)
response := answer.(*ability.HTTPResponse)

body := response.Body()
statusCode := response.StatusCode()
contentType := response.Headers()["Content-Type"]
```

## A common scenario: authenticating, then reusing a token
Log in to obtain a bearer token, then set it as the `Authorization` header for the
following requests:
```go
anActor := screenplay.ActorNamed("Wanda").WhoCan(MakeHTTPRequests())

err := anActor.AttemptsTo(SendPostRequest().To(loginURL).WithAuth(username, password))

answer, _ := last.Of(Responses()).AnsweredBy(anActor)
bearerToken := answer.(*ability.HTTPResponse).Body()

err = anActor.AttemptsTo(AddHeader("Authorization", "Bearer "+bearerToken).Secretly())
```
