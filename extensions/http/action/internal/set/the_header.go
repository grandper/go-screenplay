package set

import (
	"fmt"
	"strings"

	"github.com/grandper/go-screenplay/v2/extensions/http/action/internal/verdict"

	"github.com/grandper/go-screenplay/v2/extensions/http/action"
	impl "github.com/grandper/go-screenplay/v2/extensions/http/internal/ability"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

type headerPair struct {
	key, value string
}

// TheHeader sets a header to the actor's HTTP session:
//
//	set.TheHeader("Authorization", "Bearer "+token).Secretly()
//
// This action will remove any headers previously set in the session.
func TheHeader(key, value string) *TheHeadersAction {
	if key == "" {
		return &TheHeadersAction{err: fmt.Errorf("%w: a key cannot be empty", action.ErrInvalidHeaders)}
	}
	return &TheHeadersAction{
		headers: []headerPair{{key, value}},
		secret:  false,
	}
}

// TheHeaders sets several headers to the actor's HTTP session.
// This action will remove any headers previously set in the session.
func TheHeaders(args ...string) *TheHeadersAction {
	numArgs := len(args)
	if numArgs == 0 {
		return &TheHeadersAction{
			err: fmt.Errorf(
				"%w: set.TheHeaders should receive at least one key and its value",
				action.ErrInvalidHeaders,
			),
		}
	}
	if numArgs%2 != 0 {
		return &TheHeadersAction{
			err: fmt.Errorf("%w: set.TheHeaders should receive a list of key-value pairs", action.ErrInvalidHeaders),
		}
	}

	numHeaders := numArgs / 2 //nolint:mnd // the number is explicitly divided by 2

	headers := make([]headerPair, 0, numHeaders)
	for i := 0; i < numArgs; i += 2 {
		headers = append(headers, headerPair{args[i], args[i+1]})
	}
	return &TheHeadersAction{
		headers: headers,
		secret:  false,
	}
}

// TheHeadersAction sets headers to the actor's HTTP session.
type TheHeadersAction struct {
	headers []headerPair
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
		return fmt.Sprintf("set the headers %s", a.logHeaders())
	}
	return fmt.Sprintf("set the header %s", a.logHeaders())
}

func (a *TheHeadersAction) logHeaders() string {
	headerStrs := make([]string, 0, len(a.headers))
	if a.secret {
		for _, h := range a.headers {
			headerStrs = append(headerStrs, verdict.Describe(h.key)+" = <secret>")
		}
		return strings.Join(headerStrs, ", ")
	}
	for _, h := range a.headers {
		headerStrs = append(headerStrs, fmt.Sprintf("%s = %s", verdict.Describe(h.key), verdict.Describe(h.value)))
	}
	return strings.Join(headerStrs, ", ")
}

// PerformAs performs the task or the action as the provided actor.
func (a *TheHeadersAction) PerformAs(theActor *screenplay.Actor) error {
	if a.err != nil {
		return a.err
	}

	makeHTTPRequests := impl.TheAbilityOf(theActor)
	if err := makeHTTPRequests.ToResetHeaders(); err != nil {
		return err
	}
	for _, h := range a.headers {
		if err := makeHTTPRequests.ToAddTheHeader(h.key, h.value); err != nil {
			return err
		}
	}
	return nil
}

// TheHeadersAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*TheHeadersAction)(nil)
