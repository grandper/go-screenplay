# go-screenplay-http

This extension provides new ability, actions, questions, and resolutions
to perform HTTP requests.

## New Abilities
This extension introduces a unique capability to make api request.
To use the new capability use
```go
anActor := screeplay.ActorNamed("boby").WhoCan(MakeHTTPRequests())
```

## New Actions

### Headers
You can add of modify a header using the following actions:
```go
err := anActor.AttemptsTo(AddHeader("Content-Type", "application/json"))
err := anActor.AttemptsTo(SetHeader("Content-Type", "application/json"))
err := anActor.AttemptsTo(AddHeaders(
    "Content-Type", "application/json",
    "Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686"))
```
If you don't want the content of a header to be displayed you can make it secret:
```go
err := anActor.AttemptsTo(AddHeader("Content-Type", "application/json").Secretly())
err := anActor.AttemptsTo(AddHeader("Content-Type", "application/json").WhichShouldBeKeptSecret())
err := anActor.AttemptsTo(SetHeader("Content-Type", "application/json").Secretly())
err := anActor.AttemptsTo(SetHeader("Content-Type", "application/json").WhichShouldBeKeptSecret())
err := anActor.AttemptsTo(AddHeaders(
    "Content-Type", "application/json",
    "Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686").Secretly())
err := anActor.AttemptsTo(AddHeaders(
"Content-Type", "application/json",
"Authorization", "Bearer 84e7750a-582f-4ed7-9510-6e181d530686").WhichShouldBeKeptSecret())
```

### Sending Requests
To send request use simply need to use
```go
err := anActor.AttemptsTo(SendHTTPRequest(http.MethodGet).To("http://www.example.com").WithBody(body))
```
An HTTP request can be sent secretly, which means that the request body and headers will not be displayed in the output.
```go
err := anActor.AttemptsTo(SendHTTPRequest(http.MethodGet, "http://www.example.com").WithBody(body).Secretly())
err := anActor.AttemptsTo(SendHTTPRequest(http.MethodGet, "http://www.example.com").WithBody(body).WhichShouldBeKeptSecret())
```

For readability sake, you can also use the following shortcuts.
```go
err := anActor.AttemptsTo(SendDeleteRequest().To("http://www.example.com"))
err := anActor.AttemptsTo(SendGetRequest().To("http://www.example.com"))
err := anActor.AttemptsTo(SendHeadRequest().To("http://www.example.com"))
err := anActor.AttemptsTo(SendOptionsRequest().To("http://www.example.com"))
err := anActor.AttemptsTo(SendPatchRequest().To("http://www.example.com").WithBody(body))
err := anActor.AttemptsTo(SendPostRequest().To("http://www.example.com").WithBody(body))
err := anActor.AttemptsTo(SendPutRequest().To("http://www.example.com").WithBody(body))
```
You can use basic authentication when sending a request.
```go
err := anActor.AttemptsTo(SendGetRequest().To("http://www.example.com").WithAuth("username", "password"))
err := anActor.AttemptsTo(SendGetRequest().To("http://www.example.com").WithCredential("username", "password"))
```

## New Questions
After sending a request, an actor can ask questions about its response.

You can request information about the status code of the previous request.
```go
err := anActor.AttemptsTo(see.The(StatusCode(), is.EqualTo(200)))
err := anActor.AttemptsTo(see.The(StatusCodeOfTheLastResponse(), is.EqualTo(200)))
```

You can request information about the headers of the previous request.
```go
err := anActor.AttemptsTo(see.The(Headers(), contains.TheEntry("Content-Type", "application/json")))
err := anActor.AttemptsTo(see.The(HeadersOfTheLastResponse(), contains.TheEntry("Content-Type", "application/json")))
```

You can request information about the body of the previous request.
```go
err := anActor.AttemptsTo(see.The(Body(), contains.TheText("Hello World")))
err := anActor.AttemptsTo(see.The(BodyOfTheLastResponse(), contains.TheText("Hello World")))
```

## Setting An Authorization Header

A common scenario is to login to get a bearer token and use it for the following request.
```go
anActor := AnActor.WhoCan(MakeHTTPRequests())
anActor.AttemptsTo(SendPOSTRequest.To(loginURL).WithAuth(username, password))

bearerToken := BodyOfTheLastResponse.AnsweredBy(anActor)["token"]
anActor.AttemptsTo(AddHeader("Authorization", "Bearer " + bearerToken))
```
