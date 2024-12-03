package action

import (
	"fmt"
	"slices"
	"strings"

	"github.com/grandper/go-screenplay/extensions/http/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

// AddHeader adds n header to the actor's HTTP session.
func AddHeader(key, value string) *AddHeadersAction {
	if key == "" {
		panic("a key cannot be empty")
	}
	return &AddHeadersAction{
		headers: map[string]string{
			key: value,
		},
		secret: false,
	}
}

// AddHeaders adds several headers to the actor's HTTP session.
func AddHeaders(args ...string) *AddHeadersAction {
	numArgs := len(args)
	if numArgs == 0 {
		panic("AddHeaders should receive at least one key and its value")
	}
	if numArgs%2 != 0 {
		panic("AddHeaders should receive a list of key-value pairs")
	}

	numHeaders := numArgs / 2 //nolint:mnd // the number is explicitly divided by 2

	headers := make(map[string]string, numHeaders)
	for i := 0; i < numArgs; i += 2 {
		headers[args[i]] = args[i+1]
	}
	return &AddHeadersAction{
		headers: headers,
		secret:  false,
	}
}

// AddHeadersAction adds n header to the actor's HTTP session.
type AddHeadersAction struct {
	headers map[string]string
	secret  bool
}

// WhichShouldBeKeptSecret makes sure that the header value is not displayed in logs.
func (a *AddHeadersAction) WhichShouldBeKeptSecret() *AddHeadersAction {
	a.secret = true
	return a
}

// Secretly makes sure that the header value is not displayed in logs.
// Secretly is an alias for the method WhichShouldBeKeptSecret.
func (a *AddHeadersAction) Secretly() *AddHeadersAction {
	a.secret = true
	return a
}

// String describes the action.
func (a *AddHeadersAction) String() string {
	if len(a.headers) > 1 {
		return fmt.Sprintf("add the headers %s", a.logHeaders())
	}
	return fmt.Sprintf("add the header %s", a.logHeaders())
}

func (a *AddHeadersAction) logHeaders() string {
	headers := make([]string, 0, len(a.headers))
	if a.secret {
		for header := range a.headers {
			headers = append(headers, header+" = <secret>")
		}
		slices.Sort(headers)
		return strings.Join(headers, ", ")
	}
	for header, value := range a.headers {
		headers = append(headers, fmt.Sprintf("%s = %s", header, value))
	}
	slices.Sort(headers)
	return strings.Join(headers, ", ")
}

// PerformAs performs the task or the action as the provided actor.
func (a *AddHeadersAction) PerformAs(theActor *screenplay.Actor) error {
	makeHTTPRequests, err := screenplay.UseAbilityTo[*ability.MakeHTTPRequestsAbility]().Of(theActor)
	if err != nil {
		return err
	}
	for key, value := range a.headers {
		makeHTTPRequests.ToAddTheHeader(key, value)
	}
	return nil
}

// AddHeadersAction implements the screenplay.Action interface.
var _ screenplay.Performable = (*AddHeadersAction)(nil)
