package sentence

import (
	"errors"
	"fmt"
	"strings"

	"github.com/grandper/go-screenplay/v2/internal/verdict"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// ErrInvalidTerm is wrapped by the error returned when a sentence was built
// wrongly: a verb without a resolution, a nil resolution, a nil question, a
// nil term. It is exposed to the users as term.ErrInvalidTerm.
var ErrInvalidTerm = errors.New("invalid term")

// Verb is the narration and the polarity a clause gives its resolutions.
type Verb int

// The verbs a sentence is closed with.
const (
	Should Verb = iota
	ShouldBe
	ShouldHave
	ShouldNot
	ShouldNotBe
	ShouldNotHave
	ShouldAnyOf
)

// String describes the verb as it reads in a sentence.
func (v Verb) String() string {
	switch v {
	case Should:
		return "should"
	case ShouldBe:
		return "should be"
	case ShouldHave:
		return "should have"
	case ShouldNot:
		return "should not"
	case ShouldNotBe:
		return "should not be"
	case ShouldNotHave:
		return "should not have"
	case ShouldAnyOf:
		return "should be any of:"
	default:
		return "should"
	}
}

// method returns the name of the method a user calls to pick the verb.
func (v Verb) method() string {
	switch v {
	case Should:
		return "Should"
	case ShouldBe:
		return "ShouldBe"
	case ShouldHave:
		return "ShouldHave"
	case ShouldNot:
		return "ShouldNot"
	case ShouldNotBe:
		return "ShouldNotBe"
	case ShouldNotHave:
		return "ShouldNotHave"
	case ShouldAnyOf:
		return "ShouldAnyOf"
	default:
		return "Should"
	}
}

// asAFact returns the verb as it reads when it states a fact about a singular
// or a plural subject, and the method a user calls to pick it.
func (v Verb) asAFact(plural bool) (string, string) {
	if plural {
		switch v {
		case Should:
			return "happen to", "HappenTo"
		case ShouldBe:
			return "happen to be", "HappenToBe"
		case ShouldHave:
			return "happen to have", "HappenToHave"
		case ShouldNot:
			return "do not happen to", "DoNotHappenTo"
		case ShouldNotBe:
			return "do not happen to be", "DoNotHappenToBe"
		case ShouldNotHave:
			return "do not happen to have", "DoNotHappenToHave"
		case ShouldAnyOf:
			return "happen to be any of:", "HappenToBeAnyOf"
		default:
			return "happen to", "HappenTo"
		}
	}

	switch v {
	case Should:
		return "happens to", "HappensTo"
	case ShouldBe:
		return "happens to be", "HappensToBe"
	case ShouldHave:
		return "happens to have", "HappensToHave"
	case ShouldNot:
		return "does not happen to", "DoesNotHappenTo"
	case ShouldNotBe:
		return "does not happen to be", "DoesNotHappenToBe"
	case ShouldNotHave:
		return "does not happen to have", "DoesNotHappenToHave"
	case ShouldAnyOf:
		return "happens to be any of:", "HappensToBeAnyOf"
	default:
		return "happens to", "HappensTo"
	}
}

// negated reports whether the verb forbids its resolutions.
func (v Verb) negated() bool {
	return v == ShouldNot || v == ShouldNotBe || v == ShouldNotHave
}

// Clause is the second half of a sentence: a verb and the resolutions it
// governs, as in "should be equal to 200".
type Clause[T any] struct {
	verb        Verb
	words       string // the verb as it reads in the sentence: "should be", "happens to"
	method      string // the method the user called to pick the verb: "ShouldBe", "HappensTo"
	resolutions []screenplay.Resolution[T]
	err         error // why the clause can never be judged, when it was built wrongly
}

// NewClause binds resolutions to a verb to state an expectation: "should be
// equal to 200". A clause without a resolution, or with a nil one, can never
// be judged: it remembers why, and reports it each time it is judged.
func NewClause[T any](verb Verb, resolutions []screenplay.Resolution[T]) Clause[T] {
	return newClause(verb, verb.String(), verb.method(), resolutions)
}

// NewCondition binds resolutions to a verb to state a fact, as a condition
// does: "happens to be equal to 200". The subject of the condition is either
// singular ("the status code happens to") or plural ("the logs happen to").
func NewCondition[T any](verb Verb, plural bool, resolutions []screenplay.Resolution[T]) Clause[T] {
	words, method := verb.asAFact(plural)

	return newClause(verb, words, method, resolutions)
}

func newClause[T any](verb Verb, words, method string, resolutions []screenplay.Resolution[T]) Clause[T] {
	clause := Clause[T]{
		verb:        verb,
		words:       words,
		method:      method,
		resolutions: make([]screenplay.Resolution[T], 0, len(resolutions)),
	}

	if len(resolutions) == 0 {
		clause.err = fmt.Errorf("%w: %s needs at least one resolution", ErrInvalidTerm, method)
	}

	for position, resolution := range resolutions {
		if resolution == nil {
			clause.err = fmt.Errorf("%w: the resolution at position %d given to %s is nil",
				ErrInvalidTerm, position, method)

			continue
		}

		clause.resolutions = append(clause.resolutions, resolution)
	}

	return clause
}

// Verdict is what judging a value gives: whether it satisfies the clause, and
// when it does not, the part of the sentence that broke the rule.
type Verdict struct {
	// Met reports whether the value satisfies the clause.
	Met bool
	// Failure describes, when the clause is not met, the part of the sentence
	// that is not satisfied and why: "the status code should be less than 300:
	// got 404".
	Failure string
	// Reason is the error of the resolution that is not satisfied, when there
	// is one.
	Reason error
	// Got describes the value that was judged.
	Got string
}

// Judge inspects a value on behalf of a subject, such as "the status code". It
// fails with the error of a clause or of a resolution that was built wrongly
// and can therefore not be judged.
//
// The value satisfies the clause when every resolution holds for Should,
// ShouldBe and ShouldHave, when none does for ShouldNot, ShouldNotBe and
// ShouldNotHave, and when at least one does for ShouldAnyOf. When it does not,
// the verdict names exactly what broke the rule: the resolution that is not
// satisfied, the resolution that is satisfied although it should not be, or
// the whole list when none of them is.
func (c Clause[T]) Judge(subject string, value T) (Verdict, error) {
	if c.err != nil {
		return Verdict{}, c.err
	}

	got := verdict.Describe(value)
	unmet := Verdict{Got: got}
	held := 0

	for _, resolution := range c.resolutions {
		err := resolution.Inspect(value)

		switch {
		case errors.Is(err, screenplay.ErrInvalidResolution):
			return Verdict{}, err
		case err == nil:
			held++

			if c.verb.negated() && unmet.Failure == "" {
				unmet.Failure = fmt.Sprintf("%s %s %s: got %s", subject, c.words, resolution, got)
			}
		case !c.verb.negated() && c.verb != ShouldAnyOf && unmet.Failure == "":
			unmet.Failure = fmt.Sprintf("%s %s %s: %s", subject, c.words, resolution, err)
			unmet.Reason = err
		}
	}

	if c.verb == ShouldAnyOf && held == 0 {
		unmet.Failure = fmt.Sprintf("%s (got %s) %s", subject, got, c)
	}

	if unmet.Failure != "" {
		return unmet, nil
	}

	return Verdict{Met: true, Got: got}, nil
}

// String describes the clause as it reads after its subject: "should be equal
// to 200", "should not contain substring "panic"".
func (c Clause[T]) String() string {
	descriptions := make([]string, len(c.resolutions))
	for i, resolution := range c.resolutions {
		descriptions[i] = resolution.String()
	}

	switch {
	case c.verb == ShouldAnyOf:
		for i, description := range descriptions {
			// The copula is already in "should be any of".
			descriptions[i] = strings.TrimPrefix(description, "be ")
		}

		return c.words + " " + strings.Join(descriptions, ", ")
	case c.verb.negated():
		return c.words + " " + strings.Join(descriptions, ", nor ")
	default:
		return c.words + " " + strings.Join(descriptions, ", and ")
	}
}
