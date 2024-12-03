package ability

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

// MakeHTTPRequests enables the actor to make HTTP requests.
func MakeHTTPRequests() *MakeHTTPRequestsAbility {
	return &MakeHTTPRequestsAbility{
		headers:   map[string]string{},
		responses: []*HTTPResponse{},
	}
}

// MakeHTTPRequestsAbility is the ability to make HTTP requests.
type MakeHTTPRequestsAbility struct {
	headers   map[string]string
	responses []*HTTPResponse
}

// ToAddTheHeader adds a header to the HTTP activity.
func (mhr *MakeHTTPRequestsAbility) ToAddTheHeader(key, value string) {
	mhr.headers[key] = value
}

// ResetHeaders removes any headers previously set in the session.
func (mhr *MakeHTTPRequestsAbility) ResetHeaders() {
	mhr.headers = map[string]string{}
}

// ToRetrieveHeaders returns the headers to be applied to the next HTTP request.
func (mhr *MakeHTTPRequestsAbility) ToRetrieveHeaders() map[string]string {
	return mhr.headers
}

// ToSend sends a request. The response is stored in the ability struct.
func (mhr *MakeHTTPRequestsAbility) ToSend(method, url string, body io.Reader, credential *Credential) error {
	if !methodIsValid(method) {
		return fmt.Errorf("%s is not a valid HTTP method", method)
	}
	response, err := mhr.sendRequest(method, url, body, credential)
	if err != nil {
		return err
	}
	defer func() {
		errClose := response.Body.Close()
		if errClose != nil {
			slog.Default().ErrorContext(context.Background(), "Error closing response body", "error", errClose.Error())
		}
	}()
	httpResponse, err := NewHTTPResponseFrom(response)
	if err != nil {
		return err
	}
	mhr.responses = append(mhr.responses, httpResponse)
	return nil
}

// ToRetrieveResponses returns the responses of previous HTTP requests.
func (mhr *MakeHTTPRequestsAbility) ToRetrieveResponses() []*HTTPResponse {
	return mhr.responses
}

func methodIsValid(method string) bool {
	var validMethods = [...]string{
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
	}
	for _, validMethod := range validMethods {
		if method == validMethod {
			return true
		}
	}
	return false
}

type Credential struct {
	Username string
	Password string
}

func (mhr *MakeHTTPRequestsAbility) sendRequest(
	method, url string,
	body io.Reader,
	credential *Credential,
) (*http.Response, error) {
	request, err := http.NewRequestWithContext(context.Background(), method, url, body)
	if err != nil {
		return nil, err
	}

	if credential != nil {
		request.SetBasicAuth(credential.Username, credential.Password)
	}

	for key, value := range mhr.headers {
		request.Header.Add(key, value)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	return response, err
}

// Send sends a request. The response is stored in the ability struct.
// Send is an alias for ToSend.
func (mhr *MakeHTTPRequestsAbility) Send(method, url string, body io.Reader, credential *Credential) error {
	return mhr.ToSend(method, url, body, credential)
}

// Forget clean up the ability.
// The ability cannot be used after Forget() has been called.
// This method is used, e.g., to close connections to databases,
// deleting data, closing client cleanly.
func (mhr *MakeHTTPRequestsAbility) Forget() error {
	mhr.headers = nil
	mhr.responses = nil
	return nil
}

// String describes the ability.
func (mhr *MakeHTTPRequestsAbility) String() string {
	return "make Http requests"
}
