# go-screenplay HTTP extension

This extension lets an actor drive REST APIs. It bundles the ability to make HTTP
requests together with the actions, questions, and resolutions that go with it, so
you test an API with the same screenplay vocabulary you use everywhere else.

## The ability
This extension introduces a single ability: making HTTP requests. Give it to an actor
with `WhoCan`:
```go
anActor := screenplay.ActorNamed("Wanda").WhoCan(MakeHTTPRequests())
```
The ability keeps the headers set for the session and records every response it
receives, so later questions can ask about them.

The ability is configured fluently, and the chain reads as English: "make HTTP
requests with timing out after 5 seconds with at most 3 redirects":
```go
anActor := screenplay.ActorNamed("Wanda").WhoCan(
    MakeHTTPRequests().WithTimingOutAfter(5).Seconds().WithAtMost(3).Redirects(),
)
```
`WithTimingOutAfter` bounds every request the actor sends. Without it, a request is
bounded only by the deadline of the actor's context (see `WithContext`). A request
that outlasts the timeout fails:
```go
err := anActor.AttemptsTo(send.GetRequest().To("https://slow.example.com/report"))
// Get "https://slow.example.com/report": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
```
The timeout is closed by its unit: `Milliseconds`, `Seconds` or `Minutes`, each with
its singular form (`Millisecond`, `Second`, `Minute`) so that one unit reads as English:
```go
anActor := screenplay.ActorNamed("Wanda").WhoCan(MakeHTTPRequests().WithTimingOutAfter(500).Milliseconds())
anActor := screenplay.ActorNamed("Wanda").WhoCan(MakeHTTPRequests().WithTimingOutAfter(1).Minute())
```
`WithAtMost(n).Redirects()` bounds the number of consecutive redirects followed
before a request fails. Without it, ten redirects are followed, as `net/http` does
by default; zero makes every redirect fail:
```go
err := anActor.AttemptsTo(send.GetRequest().To("https://www.example.com/loop"))
// Get "https://www.example.com/loop": stopped after 3 redirects
```

The settings are the only methods a test can call on the ability. Its operations live
in `extensions/http/internal/ability`, so only an action or a question of the extension
uses them, in one statement written in `PerformAs` or `AnsweredBy` that reads "the
ability of the actor to add the header":
```go
err := impl.TheAbilityOf(theActor).ToAddTheHeader("Accept", "application/json")
```
`TheAbilityOf` never fails: when the actor cannot make HTTP requests, the operation
reports `screenplay.ErrActorMissingAbility`.

When the actor exits, the ability forgets the session: the headers and the recorded
responses are cleared.
```go
err := anActor.Exit()
```

## Actions

### Headers
Add a header (keeping any header already set), or set the headers (replacing the
whole set) for every subsequent request of the session:
```go
err := anActor.AttemptsTo(add.TheHeader("Content-Type", "application/json"))
// Wanda attempts to add the header "Content-Type" = "application/json"
err := anActor.AttemptsTo(add.TheHeaders(
    "Content-Type", "application/json",
    "Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686"))
// Wanda attempts to add the headers "Authorization" = "Bearer 84e7750a-582f-4ed7-9510-6e181d530686", "Content-Type" = "application/json"
err := anActor.AttemptsTo(set.TheHeader("Content-Type", "application/json"))
// Wanda attempts to set the header "Content-Type" = "application/json"
err := anActor.AttemptsTo(set.TheHeaders(
    "Content-Type", "application/json",
    "Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686"))
// Wanda attempts to set the headers "Content-Type" = "application/json", "Authorization" = "Bearer 84e7750a-582f-4ed7-9510-6e181d530686"
```
`add.TheHeaders` narrates its headers sorted by key; `set.TheHeaders` narrates them in
the order they are written.
`add.TheHeaders` and `set.TheHeaders` take an even number of arguments, read as
`key, value, key, value, ...`.

Nothing panics when the headers are written wrongly. An empty key, no header at all, or a
key without its value is reported when the actor performs the action, with an error wrapping
`ErrInvalidHeaders`:
```go
err := anActor.AttemptsTo(add.TheHeaders("Content-Type"))
// invalid headers: add.TheHeaders should receive a list of key-value pairs
err := anActor.AttemptsTo(add.TheHeaders())
// invalid headers: add.TheHeaders should receive at least one key and its value
err := anActor.AttemptsTo(set.TheHeader("", "application/json"))
// invalid headers: a key cannot be empty
```

Headers are sent with every request of the actor. When a server redirects to another
host, the headers added through the ability are not forwarded to it.

If a header holds a secret you do not want in the narration, mark the action secret
with `Secretly` (or its longer alias `WhichShouldBeKeptSecret`):
```go
err := anActor.AttemptsTo(add.TheHeader("Authorization", "Bearer "+token).Secretly())
err := anActor.AttemptsTo(set.TheHeader("Authorization", "Bearer "+token).WhichShouldBeKeptSecret())
// Wanda attempts to set the header "Authorization" = <secret>
```

### Sending requests
Send a request with `send.RequestWithMethod`, choosing the method, the URL (`To`), and
optionally a body (`WithBody`):
```go
err := anActor.AttemptsTo(send.RequestWithMethod(http.MethodGet).To("http://www.example.com"))
// Wanda attempts to send a GET request to "http://www.example.com"
err := anActor.AttemptsTo(send.RequestWithMethod(http.MethodPost).To("http://www.example.com").WithBody(strings.NewReader(`{"name":"Wanda"}`)))
// Wanda attempts to send a POST request to "http://www.example.com" with body "{\"name\":\"Wanda\"}"
```
The method is one of `DELETE`, `GET`, `HEAD`, `OPTIONS`, `PATCH`, `POST` and `PUT`. Any
other method is reported when the actor performs the action, with an error wrapping
`ErrInvalidHTTPMethod`:
```go
err := anActor.AttemptsTo(send.RequestWithMethod("FETCH").To("http://www.example.com"))
// "FETCH" is not a valid HTTP method
```
For readability there is a shortcut per method:
```go
err := anActor.AttemptsTo(send.DeleteRequest().To("http://www.example.com"))
err := anActor.AttemptsTo(send.GetRequest().To("http://www.example.com"))
err := anActor.AttemptsTo(send.HeadRequest().To("http://www.example.com"))
err := anActor.AttemptsTo(send.OptionsRequest().To("http://www.example.com"))
err := anActor.AttemptsTo(send.PatchRequest().To("http://www.example.com").WithBody(body))
err := anActor.AttemptsTo(send.PostRequest().To("http://www.example.com").WithBody(body))
err := anActor.AttemptsTo(send.PutRequest().To("http://www.example.com").WithBody(body))
```
`WithBody` takes an `io.Reader`. It buffers the reader so the request can be
described in the narration and sent (and re-sent) without the body being consumed.
A body larger than 64 MiB is reported when the actor performs the action:
```go
err := anActor.AttemptsTo(send.PostRequest().To("http://www.example.com/upload").WithBody(hugeArchive))
// failed to read the request body: the request body exceeds 64 MiB
```
So is a body that cannot be read: nothing is sent, the error of the reader is reported,
and the narration says so instead of quoting a body that never existed:
```go
err := anActor.AttemptsTo(send.PostRequest().To("http://www.example.com/upload").WithBody(brokenArchive))
// Wanda attempts to send a POST request to "http://www.example.com/upload" with a body that could not be read
// failed to read the request body: the disk is gone
```

Use basic authentication with `WithCredential` (or its alias `WithAuth`):
```go
err := anActor.AttemptsTo(send.GetRequest().To("http://www.example.com").WithCredential("username", "password"))
err := anActor.AttemptsTo(send.GetRequest().To("http://www.example.com").WithAuth("username", "password"))
```
A request can be sent secretly, so neither its URL nor its body appear in the
narration. When a secret request fails, its error keeps the cause but drops the
URL, so a token carried in the query string never reaches the logs:
```go
err := anActor.AttemptsTo(send.PostRequest().To("http://www.example.com").WithBody(body).Secretly())
err := anActor.AttemptsTo(send.PostRequest().To("http://www.example.com").WithBody(body).WhichShouldBeKeptSecret())
// narrated: Wanda attempts to send a secret HTTP request
// on failure: Wanda attempts to send a secret HTTP request: failed to send a secret HTTP request: dial tcp 127.0.0.1:1: connect: connection refused
```

## Questions

### The responses
After sending one or more requests, `Responses()` answers with every
`*ability.HTTPResponse` the actor received, in the order the requests were sent.
Because it answers with a slice, it composes with the core `first`/`last`/`number`
questions to focus on a single response:
```go
err := anActor.Should(see.The(number.Of(Responses())).Should(be.EqualTo(2)))
```
A response body is read whole when the response arrives. A body larger than 64 MiB
makes the request fail, so a large download never exhausts the memory of the test:
```go
err := anActor.AttemptsTo(send.GetRequest().To("http://www.example.com/backup.tar"))
// failed to create a new HTTPResponse: the response body exceeds 64 MiB
```

### The status code
`status.CodeOf` wraps a question that answers with an `*HTTPResponse` and answers
with its status code, so you assert on it like any number:
```go
err := anActor.Should(see.The(status.CodeOf(last.Of(Responses()))).Should(be.EqualTo(200)))
```
When the wrapped question answers with a nil response, the question fails with an
error wrapping `status.ErrResponseIsNil`, followed by the description of the wrapped
question, here a question described as "cached response":
```go
err := anActor.Should(see.The(status.CodeOf(cachedResponse)).Should(be.EqualTo(200)))
// the answer of the question is a nil HTTP response: cached response
```

### Reading a response directly
An `*ability.HTTPResponse` exposes its parts through getters — `Body()` (a
`string`), `Headers()` (a `map[string]string`), and `StatusCode()` (an `int`):
```go
response, _ := anActor.AsksFor(last.Of(Responses())) // response is a *ability.HTTPResponse

body := response.Body()
statusCode := response.StatusCode()
contentType := response.Headers()["Content-Type"]
```
A header the server sent several times holds all its values, joined with a comma.
`Headers()` is keyed by the canonical name of each header (`Content-Type`, not
`content-type`).

## A common scenario: authenticating, then reusing a token
Log in to obtain a bearer token, then set it as the `Authorization` header for the
following requests:
```go
anActor := screenplay.ActorNamed("Wanda").WhoCan(MakeHTTPRequests())

err := anActor.AttemptsTo(send.PostRequest().To(loginURL).WithAuth(username, password))

response, _ := anActor.AsksFor(last.Of(Responses()))
bearerToken := response.Body()

err = anActor.AttemptsTo(add.TheHeader("Authorization", "Bearer "+bearerToken).Secretly())
```
