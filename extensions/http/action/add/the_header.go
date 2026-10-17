package add

import (
	"fmt"
	"slices"
	"strings"

	"github.com/grandper/go-screenplay/v2/internal/verdict"

	"github.com/grandper/go-screenplay/v2/extensions/http/ability"
	"github.com/grandper/go-screenplay/v2/extensions/http/action"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// TheHeader adds a header to the actor's HTTP session:
//
//	add.TheHeader("Accept", "application/json")
func TheHeader(key, value string) *TheHeadersAction {
	if key == "" {
		return &TheHeadersAction{err: fmt.Errorf("%w: a key cannot be empty", action.ErrInvalidHeaders)}
	}
	return &TheHeadersAction{
		headers: map[string]string{
			key: value,
		},
		secret: false,
	}
}

// TheHeaders adds several headers to the actor's HTTP session.
func TheHeaders(args ...string) *TheHeadersAction {
	numArgs := len(args)
	if numArgs == 0 {
		return &TheHeadersAction{
			err: fmt.Errorf("%w: add.TheHeaders should receive at least one key and its value", action.ErrInvalidHeaders),
		}
	}
	if numArgs%2 != 0 {
		return &TheHeadersAction{
			err: fmt.Errorf("%w: add.TheHeaders should receive a list of key-value pairs", action.ErrInvalidHeaders),
		}
	}

	numHeaders := numArgs / 2 //nolint:mnd // the number is explicitly divided by 2

	headers := make(map[string]string, numHeaders)
	for i := 0; i < numArgs; i += 2 {
		headers[args[i]] = args[i+1]
	}
	return &TheHeadersAction{
		headers: headers,
		secret:  false,
	}
}

// TheHeadersAction adds headers to the actor's HTTP session.
type TheHeadersAction struct {
	headers map[string]string
	secret  bool
	err     error // why the action cannot be performed, when it was built wrongly
}

// WhichShouldBeKeptSecret makes sure that the header value is not displayed in logs.
func (a *TheHeadersAction) WhichShouldBeKeptSecret() *TheHeadersAction {
	a.secret = true
	return a
}

// Secretly makes sure that the header value is not displayed in logs.
// Secretly is an alias for the method WhichShouldBeKeptSecret.
func (a *TheHeadersAction) Secretly() *TheHeadersAction {
	a.secret = true
	return a
}

// String describes the action.
func (a *TheHeadersAction) String() string {
	if len(a.headers) > 1 {
		return fmt.Sprintf("add the headers %s", a.logHeaders())
	}
	return fmt.Sprintf("add the header %s", a.logHeaders())
}

func (a *TheHeadersAction) logHeaders() string {
	headers := make([]string, 0, len(a.headers))
	if a.secret {
		for header := range a.headers {
			headers = append(headers, verdict.Describe(header)+" = <secret>")
		}
		slices.Sort(headers)
		return strings.Join(headers, ", ")
	}
	for header, value := range a.headers {
		headers = append(headers, fmt.Sprintf("%s = %s", verdict.Describe(header), verdict.Describe(value)))
	}
	slices.Sort(headers)
	return strings.Join(headers, ", ")
}

// PerformAs performs the task or the action as the provided actor.
func (a *TheHeadersAction) PerformAs(theActor *screenplay.Actor) error {
	if a.err != nil {
		return a.err
	}

	makeHTTPRequests := screenplay.UseAnAbilityOf(theActor).That(ability.MakeHTTPRequests())
	for key, value := range a.headers {
		if err := makeHTTPRequests.ToAddTheHeader(key, value); err != nil {
			return err
		}
	}
	return nil
}

// TheHeadersAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheHeadersAction)(nil)
