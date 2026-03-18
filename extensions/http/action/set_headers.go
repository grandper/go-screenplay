package action

import (
	"fmt"
	"strings"

	"github.com/grandper/go-screenplay/extensions/http/ability"
	"github.com/grandper/go-screenplay/screenplay"
)

type headerPair struct {
	key, value string
}

// SetHeader sets a header to the actor's HTTP session.
// This actions will remove any headers previously set in the session.
func SetHeader(key, value string) *SetHeadersAction {
	if key == "" {
		panic("a key cannot be empty")
	}
	return &SetHeadersAction{
		headers: []headerPair{{key, value}},
		secret:  false,
	}
}

// SetHeaders sets several headers to the actor's HTTP session.
// This actions will remove any headers previously set in the session.
func SetHeaders(args ...string) *SetHeadersAction {
	numArgs := len(args)
	if numArgs == 0 {
		panic("SetHeaders should receive at least one key and its value")
	}
	if numArgs%2 != 0 {
		panic("SetHeaders should receive a list of key-value pairs")
	}

	numHeaders := numArgs / 2 //nolint:mnd // the number is explicitly divided by 2

	headers := make([]headerPair, 0, numHeaders)
	for i := 0; i < numArgs; i += 2 {
		headers = append(headers, headerPair{args[i], args[i+1]})
	}
	return &SetHeadersAction{
		headers: headers,
		secret:  false,
	}
}

// SetHeadersAction sets a header to the actor's HTTP session.
type SetHeadersAction struct {
	headers []headerPair
	secret  bool
}

// WhichShouldBeKeptSecret makes sure that the header value is not displayed in logs.
func (a *SetHeadersAction) WhichShouldBeKeptSecret() *SetHeadersAction {
	a.secret = true
	return a
}

// Secretly makes sure that the header value is not displayed in logs.
// Secretly is an alias for the method WhichShouldBeKeptSecret.
func (a *SetHeadersAction) Secretly() *SetHeadersAction {
	a.secret = true
	return a
}

// String describes the action.
func (a *SetHeadersAction) String() string {
	if len(a.headers) > 1 {
		return fmt.Sprintf("set the headers %s", a.logHeaders())
	}
	return fmt.Sprintf("set the header %s", a.logHeaders())
}

func (a *SetHeadersAction) logHeaders() string {
	headerStrs := make([]string, 0, len(a.headers))
	if a.secret {
		for _, h := range a.headers {
			headerStrs = append(headerStrs, h.key+" = <secret>")
		}
		return strings.Join(headerStrs, ", ")
	}
	for _, h := range a.headers {
		headerStrs = append(headerStrs, fmt.Sprintf("%s = %s", h.key, h.value))
	}
	return strings.Join(headerStrs, ", ")
}

// PerformAs performs the task or the action as the provided actor.
func (a *SetHeadersAction) PerformAs(theActor *screenplay.Actor) error {
	makeHTTPRequests, err := screenplay.UseAbilityTo[*ability.MakeHTTPRequestsAbility]().Of(theActor)
	if err != nil {
		return err
	}
	makeHTTPRequests.ResetHeaders()
	for _, h := range a.headers {
		makeHTTPRequests.ToAddTheHeader(h.key, h.value)
	}
	return nil
}

// SetHeadersAction implements the screenplay.Action interface.
var _ screenplay.Performable = (*SetHeadersAction)(nil)
