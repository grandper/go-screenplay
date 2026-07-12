package action

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/grandper/go-screenplay/extensions/http/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// SendHTTPRequest sends an HTTP request.
func SendHTTPRequest(method string) *SendHTTPRequestAction {
	return &SendHTTPRequestAction{
		method: method,
		url:    "127.0.0.1",
		secret: false,
	}
}

// SendDeleteRequest sends a DELETE request.
func SendDeleteRequest() *SendHTTPRequestAction {
	return SendHTTPRequest(http.MethodDelete)
}

// SendGetRequest sends a GET request.
func SendGetRequest() *SendHTTPRequestAction {
	return SendHTTPRequest(http.MethodGet)
}

// SendHeadRequest sends a HEAD request.
func SendHeadRequest() *SendHTTPRequestAction {
	return SendHTTPRequest(http.MethodHead)
}

// SendOptionsRequest sends an OPTIONS request.
func SendOptionsRequest() *SendHTTPRequestAction {
	return SendHTTPRequest(http.MethodOptions)
}

// SendPatchRequest sends a PATCH request.
func SendPatchRequest() *SendHTTPRequestAction {
	return SendHTTPRequest(http.MethodPatch)
}

// SendPostRequest sends a POST request.
func SendPostRequest() *SendHTTPRequestAction {
	return SendHTTPRequest(http.MethodPost)
}

// SendPutRequest sends a PUT request.
func SendPutRequest() *SendHTTPRequestAction {
	return SendHTTPRequest(http.MethodPut)
}

// SendHTTPRequestAction sends an HTTP request.
type SendHTTPRequestAction struct {
	method     string
	url        string
	body       []byte
	hasBody    bool
	bodyErr    error
	credential *ability.Credential
	secret     bool
}

// To sets the URL to which the request should be sent.
func (a *SendHTTPRequestAction) To(url string) *SendHTTPRequestAction {
	a.url = url
	return a
}

// WithBody sets the body of the request. The reader is buffered so that the
// body can be described (String) and sent (PerformAs) without being consumed,
// and so the action can be performed more than once.
func (a *SendHTTPRequestAction) WithBody(body io.Reader) *SendHTTPRequestAction {
	if body == nil {
		return a
	}

	a.hasBody = true
	a.body, a.bodyErr = io.ReadAll(body)

	return a
}

// WithCredential sets the basic authentication credential to use for the request.
func (a *SendHTTPRequestAction) WithCredential(username, password string) *SendHTTPRequestAction {
	a.credential = &ability.Credential{
		Username: username,
		Password: password,
	}
	return a
}

// WithAuth is an alias for WithCredential.
func (a *SendHTTPRequestAction) WithAuth(username, password string) *SendHTTPRequestAction {
	return a.WithCredential(username, password)
}

// WhichShouldBeKeptSecret makes sure that the request is not displayed in logs.
func (a *SendHTTPRequestAction) WhichShouldBeKeptSecret() *SendHTTPRequestAction {
	a.secret = true
	return a
}

// Secretly makes sure that the request is not displayed in logs.
// Secretly is an alias for the method WhichShouldBeKeptSecret.
func (a *SendHTTPRequestAction) Secretly() *SendHTTPRequestAction {
	a.secret = true
	return a
}

// String describes the action.
func (a *SendHTTPRequestAction) String() string {
	if a.secret {
		return "send a secret HTTP request"
	}
	if a.hasBody {
		body := string(a.body)
		if a.bodyErr != nil {
			body = "failed to read the body"
		}
		return fmt.Sprintf("send a %s request to %s with body '%s'", a.method, a.url, body)
	}
	return fmt.Sprintf("send a %s request to %s", a.method, a.url)
}

// PerformAs performs the task or the action as the provided actor.
func (a *SendHTTPRequestAction) PerformAs(theActor *screenplay.Actor) error {
	if a.bodyErr != nil {
		return fmt.Errorf("failed to read the request body: %w", a.bodyErr)
	}

	makeHTTPRequests, err := screenplay.UseAbilityTo[*ability.MakeHTTPRequestsAbility]().Of(theActor)
	if err != nil {
		return err
	}

	var body io.Reader
	if a.hasBody {
		body = bytes.NewReader(a.body)
	}

	return makeHTTPRequests.ToSend(a.method, a.url, body, a.credential)
}

// SendHTTPRequestAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*SendHTTPRequestAction)(nil)
