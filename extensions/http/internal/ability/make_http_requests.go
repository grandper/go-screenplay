package ability

import (
	"context"
	"fmt"
	"io"
	"maps"
	"net/http"
	"slices"
	"sync"

	"github.com/grandper/go-screenplay/v2/screenplay"
	"github.com/grandper/go-screenplay/v2/timing"
)

// defaultMaxRedirects is the number of consecutive redirects followed before a
// request fails when WithAtMost is not used, the same limit as the default
// policy of net/http.
const defaultMaxRedirects = 10

// MakeHTTPRequests enables the actor to make HTTP requests. The ability is
// configured fluently, reading "make HTTP requests with ...":
//
//	MakeHTTPRequests().WithTimingOutAfter(5).Seconds().WithAtMost(3).Redirects()
func MakeHTTPRequests() *MakeHTTPRequestsAbility {
	mhr := &MakeHTTPRequestsAbility{
		headers:      map[string]string{},
		responses:    []*HTTPResponse{},
		maxRedirects: defaultMaxRedirects,
	}
	mhr.client = &http.Client{
		CheckRedirect: mhr.checkRedirect,
	}
	return mhr
}

// WithTimingOutAfter bounds every request the ability sends; it reads "with
// timing out after 5 seconds". Without it, a request is bounded only by the
// deadline of the actor's context.
func (mhr *MakeHTTPRequestsAbility) WithTimingOutAfter(amount int) *timing.DurationBuilder[MakeHTTPRequestsAbility] {
	return timing.NewDurationBuilder(mhr, &mhr.client.Timeout).For(amount)
}

// WithAtMost bounds the number of consecutive redirects followed before a
// request fails; it reads "with at most 3 redirects". Without it, ten
// redirects are followed, as the default policy of net/http does.
func (mhr *MakeHTTPRequestsAbility) WithAtMost(amount int) *AtMostBuilder {
	return &AtMostBuilder{ability: mhr, amount: amount}
}

// AtMostBuilder holds the amount given to WithAtMost and waits for the unit that
// closes the setting.
type AtMostBuilder struct {
	ability *MakeHTTPRequestsAbility
	amount  int
}

// Redirects makes the amount the number of consecutive redirects followed
// before a request fails. Zero makes every redirect fail.
func (b *AtMostBuilder) Redirects() *MakeHTTPRequestsAbility {
	b.ability.maxRedirects = b.amount

	return b.ability
}

// checkRedirect is the redirect policy of the ability's client. Redirects are
// followed, up to the limit set with WithAtMost, but when a redirect leaves the
// host the request was sent to, every header added through the ability is
// removed from the redirected request, so a session header never reaches
// another host.
func (mhr *MakeHTTPRequestsAbility) checkRedirect(request *http.Request, via []*http.Request) error {
	if len(via) >= mhr.maxRedirects {
		return fmt.Errorf("stopped after %d redirects", mhr.maxRedirects)
	}
	if request.URL.Host != via[0].URL.Host {
		mhr.mutex.RLock()
		defer mhr.mutex.RUnlock()

		for key := range mhr.headers {
			request.Header.Del(key)
		}
	}
	return nil
}

// MakeHTTPRequestsAbility is the ability to make HTTP requests.
// It is safe for concurrent use: actions running in parallel (see
// proceed.ConcurrentlyTo) can add headers and send requests without racing.
// Requests are sent through the ability's own client, never the process-wide
// http.DefaultClient, and carry the actor's context so that a deadline or a
// cancellation set on the actor stops a request in flight.
type MakeHTTPRequestsAbility struct {
	client       *http.Client
	maxRedirects int
	mutex        sync.RWMutex
	headers      map[string]string
	responses    []*HTTPResponse
}

// TheAbilityOf returns the operations of the actor's ability to make HTTP
// requests; it is written in PerformAs or AnsweredBy and reads "the ability of
// the actor to send":
//
//	err := ability.TheAbilityOf(theActor).ToSend(ctx, method, url, body, nil)
//
// It never fails. When the actor is nil or does not have the ability, every
// operation fails with screenplay.ErrActorMissingAbility.
func TheAbilityOf(actor *screenplay.Actor) *Operations {
	return &Operations{ability: screenplay.UseAnAbilityOf(actor).That(MakeHTTPRequests())}
}

// Operations holds the operations of the actor's ability to make HTTP requests.
// They are not methods of MakeHTTPRequestsAbility, whose method set holds only
// its settings, so a test cannot send a request without an action.
type Operations struct {
	ability *MakeHTTPRequestsAbility
}

// ToAddTheHeader adds a header to the HTTP activity. When the actor does not
// have the ability, it fails with screenplay.ErrActorMissingAbility.
func (o *Operations) ToAddTheHeader(key, value string) error {
	mhr := o.ability
	if mhr == nil {
		return screenplay.ErrActorMissingAbility
	}

	mhr.mutex.Lock()
	defer mhr.mutex.Unlock()

	// After Forget the map is nil; a branch still running then must not panic.
	if mhr.headers == nil {
		mhr.headers = map[string]string{}
	}
	mhr.headers[key] = value

	return nil
}

// ToResetHeaders removes any headers previously set in the session. When the
// actor does not have the ability, it fails with
// screenplay.ErrActorMissingAbility.
func (o *Operations) ToResetHeaders() error {
	mhr := o.ability
	if mhr == nil {
		return screenplay.ErrActorMissingAbility
	}

	mhr.mutex.Lock()
	defer mhr.mutex.Unlock()

	mhr.headers = map[string]string{}

	return nil
}

// ToRetrieveHeaders returns a copy of the headers to be applied to the next
// HTTP request. When the actor does not have the ability, it fails with
// screenplay.ErrActorMissingAbility.
func (o *Operations) ToRetrieveHeaders() (map[string]string, error) {
	mhr := o.ability
	if mhr == nil {
		return nil, screenplay.ErrActorMissingAbility
	}

	mhr.mutex.RLock()
	defer mhr.mutex.RUnlock()

	return maps.Clone(mhr.headers), nil
}

// ToSend sends a request with the given context, which is the actor's context
// when called from an action. The response is stored in the ability struct.
// When the actor does not have the ability, it fails with
// screenplay.ErrActorMissingAbility.
func (o *Operations) ToSend(
	ctx context.Context,
	method, url string,
	body io.Reader,
	credential *Credential,
) error {
	mhr := o.ability
	if mhr == nil {
		return screenplay.ErrActorMissingAbility
	}
	if !methodIsValid(method) {
		return fmt.Errorf("%q %w", method, ErrInvalidHTTPMethod)
	}
	response, err := o.sendRequest(ctx, method, url, body, credential)
	if err != nil {
		return err
	}
	httpResponse, err := NewHTTPResponseFrom(response)
	if errClose := response.Body.Close(); errClose != nil && err == nil {
		err = fmt.Errorf("failed to close the response body: %w", errClose)
	}
	if err != nil {
		return err
	}

	mhr.mutex.Lock()
	defer mhr.mutex.Unlock()

	mhr.responses = append(mhr.responses, httpResponse)
	return nil
}

// ToRetrieveResponses returns a copy of the responses of previous HTTP
// requests. When the actor does not have the ability, it fails with
// screenplay.ErrActorMissingAbility.
func (o *Operations) ToRetrieveResponses() ([]*HTTPResponse, error) {
	mhr := o.ability
	if mhr == nil {
		return nil, screenplay.ErrActorMissingAbility
	}

	mhr.mutex.RLock()
	defer mhr.mutex.RUnlock()

	return slices.Clone(mhr.responses), nil
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

// Credential holds the username and password of a basic authentication credential.
type Credential struct {
	Username string
	Password string
}

func (o *Operations) sendRequest(
	ctx context.Context,
	method, url string,
	body io.Reader,
	credential *Credential,
) (*http.Response, error) {
	mhr := o.ability
	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	if credential != nil {
		request.SetBasicAuth(credential.Username, credential.Password)
	}

	// The headers are snapshotted under the read lock, so the network call runs
	// without holding the lock and a concurrent header change never races the
	// iteration.
	headers, err := o.ToRetrieveHeaders()
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	response, err := mhr.client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, err
}

// Forget clears the headers and the responses of the session. On a nil
// receiver, the actor does not have the ability, it fails with
// screenplay.ErrActorMissingAbility.
func (mhr *MakeHTTPRequestsAbility) Forget() error {
	if mhr == nil {
		return screenplay.ErrActorMissingAbility
	}

	mhr.mutex.Lock()
	defer mhr.mutex.Unlock()

	mhr.headers = nil
	mhr.responses = nil
	return nil
}

// MakeHTTPRequestsAbility implements the screenplay.Ability interface.
var _ screenplay.Ability = (*MakeHTTPRequestsAbility)(nil)

// String describes the ability.
func (mhr *MakeHTTPRequestsAbility) String() string {
	return "make HTTP requests"
}
