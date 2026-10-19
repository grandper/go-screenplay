package proceed

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unicode"

	"github.com/grandper/go-screenplay/v2/action/internal/condition"
	"github.com/grandper/go-screenplay/v2/internal/sentence"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// RepeatedlyTo starts an action that performs a body again and again. Close it
// with the clause that says how many passes it makes:
//
//	proceed.RepeatedlyTo(TakeAMessage()).WhileThe(queueLength).HappensTo(be.GreaterThan(0))
//	proceed.RepeatedlyTo(PollTheJob()).UntilThe(jobStatus).HappensTo(be.EqualTo("done"))
//	proceed.RepeatedlyTo(AddTheBookToTheCart()).Exactly(3).Times()
//
// The body is an action, or a function of the index of the pass that returns
// an action, as for index := range 3 does in Go; the index starts at 0:
//
//	proceed.RepeatedlyTo(UploadThePart).Exactly(3).Times()
//
// Any other body builds, and fails with screenplay.ErrInvalidResolution when
// the loop is performed.
func RepeatedlyTo[B any](body B) *RepeatedlyToBuilder[B] {
	return &RepeatedlyToBuilder[B]{body: body}
}

// RepeatedlyToBuilder holds the body of a proceed.RepeatedlyTo action, and
// waits for the clause that says how many passes it makes.
type RepeatedlyToBuilder[B any] struct {
	body B
}

// AtLeastOnce performs the body once before the condition is first tested:
//
//	proceed.RepeatedlyTo(TakeAMessage()).AtLeastOnce().WhileThe(queueLength).HappensTo(be.GreaterThan(0))
func (b *RepeatedlyToBuilder[B]) AtLeastOnce() *RepeatedlyToAtLeastOnceBuilder[B] {
	return &RepeatedlyToAtLeastOnceBuilder[B]{body: b.body}
}

// WhileThe repeats the body while the answer to a question happens to be as
// stated. The condition is tested before each pass. Close it with a verb and
// one or more resolutions:
//
//	proceed.RepeatedlyTo(TakeAMessage()).WhileThe(queueLength).HappensTo(be.GreaterThan(0))
func (b *RepeatedlyToBuilder[B]) WhileThe[T any](question screenplay.Question[T]) *RepeatedlyToConditionBuilder[B, T] {
	return &RepeatedlyToConditionBuilder[B, T]{body: b.body, question: question}
}

// UntilThe repeats the body until the answer to a question happens to be as
// stated. The condition is tested before each pass. Close it with a verb and
// one or more resolutions:
//
//	proceed.RepeatedlyTo(PollTheJob()).UntilThe(jobStatus).HappensTo(be.EqualTo("done"))
func (b *RepeatedlyToBuilder[B]) UntilThe[T any](question screenplay.Question[T]) *RepeatedlyToConditionBuilder[B, T] {
	return &RepeatedlyToConditionBuilder[B, T]{body: b.body, question: question, until: true}
}

// Exactly sets the number of passes. Close it with Times:
//
//	proceed.RepeatedlyTo(AddTheBookToTheCart()).Exactly(3).Times()
func (b *RepeatedlyToBuilder[B]) Exactly(count int) *RepeatedlyToCountBuilder[B] {
	return &RepeatedlyToCountBuilder[B]{body: b.body, count: count}
}

// RepeatedlyToAtLeastOnceBuilder holds the body of a proceed.RepeatedlyTo
// action that performs it once before testing its condition, and waits for
// the condition.
type RepeatedlyToAtLeastOnceBuilder[B any] struct {
	body B
}

// WhileThe repeats the body while the answer to a question happens to be as
// stated. The condition is tested after the first pass, then before each
// other pass.
func (b *RepeatedlyToAtLeastOnceBuilder[B]) WhileThe[T any](
	question screenplay.Question[T],
) *RepeatedlyToConditionBuilder[B, T] {
	return &RepeatedlyToConditionBuilder[B, T]{body: b.body, question: question, atLeastOnce: true}
}

// UntilThe repeats the body until the answer to a question happens to be as
// stated. The condition is tested after the first pass, then before each
// other pass.
func (b *RepeatedlyToAtLeastOnceBuilder[B]) UntilThe[T any](
	question screenplay.Question[T],
) *RepeatedlyToConditionBuilder[B, T] {
	return &RepeatedlyToConditionBuilder[B, T]{body: b.body, question: question, until: true, atLeastOnce: true}
}

// RepeatedlyToConditionBuilder builds the condition of a proceed.RepeatedlyTo
// action: it holds the question the loop tests, and waits for the verb that
// closes the condition.
type RepeatedlyToConditionBuilder[B, T any] struct {
	body        B
	question    screenplay.Question[T]
	until       bool
	atLeastOnce bool
}

// HappensTo tests the answer against every resolution. It reads "happens to ...".
func (b *RepeatedlyToConditionBuilder[B, T]) HappensTo(resolutions ...screenplay.Resolution[T]) screenplay.Performable {
	return b.repeats(sentence.NewCondition(sentence.Should, false, resolutions))
}

// HappensToBe tests the answer against every resolution. It reads "happens to be ...".
func (b *RepeatedlyToConditionBuilder[B, T]) HappensToBe(
	resolutions ...screenplay.Resolution[T],
) screenplay.Performable {
	return b.repeats(sentence.NewCondition(sentence.ShouldBe, false, resolutions))
}

// HappensToHave tests the answer against every resolution. It reads "happens to have ...".
func (b *RepeatedlyToConditionBuilder[B, T]) HappensToHave(
	resolutions ...screenplay.Resolution[T],
) screenplay.Performable {
	return b.repeats(sentence.NewCondition(sentence.ShouldHave, false, resolutions))
}

// DoesNotHappenTo tests that the answer satisfies none of the resolutions. It reads "does not happen to ...".
func (b *RepeatedlyToConditionBuilder[B, T]) DoesNotHappenTo(
	resolutions ...screenplay.Resolution[T],
) screenplay.Performable {
	return b.repeats(sentence.NewCondition(sentence.ShouldNot, false, resolutions))
}

// DoesNotHappenToBe tests that the answer satisfies none of the resolutions. It reads "does not happen to be ...".
func (b *RepeatedlyToConditionBuilder[B, T]) DoesNotHappenToBe(
	resolutions ...screenplay.Resolution[T],
) screenplay.Performable {
	return b.repeats(sentence.NewCondition(sentence.ShouldNotBe, false, resolutions))
}

// DoesNotHappenToHave tests that the answer satisfies none of the resolutions. It reads "does not happen to have ...".
func (b *RepeatedlyToConditionBuilder[B, T]) DoesNotHappenToHave(
	resolutions ...screenplay.Resolution[T],
) screenplay.Performable {
	return b.repeats(sentence.NewCondition(sentence.ShouldNotHave, false, resolutions))
}

// HappensToBeAnyOf tests that the answer satisfies at least one resolution. It reads "happens to be any of: ...".
func (b *RepeatedlyToConditionBuilder[B, T]) HappensToBeAnyOf(
	resolutions ...screenplay.Resolution[T],
) screenplay.Performable {
	return b.repeats(sentence.NewCondition(sentence.ShouldAnyOf, false, resolutions))
}

// HappenTo tests the answer against every resolution. It is the plural form: it reads "happen to ...".
func (b *RepeatedlyToConditionBuilder[B, T]) HappenTo(resolutions ...screenplay.Resolution[T]) screenplay.Performable {
	return b.repeats(sentence.NewCondition(sentence.Should, true, resolutions))
}

// HappenToBe tests the answer against every resolution. It is the plural form: it reads "happen to be ...".
func (b *RepeatedlyToConditionBuilder[B, T]) HappenToBe(
	resolutions ...screenplay.Resolution[T],
) screenplay.Performable {
	return b.repeats(sentence.NewCondition(sentence.ShouldBe, true, resolutions))
}

// HappenToHave tests the answer against every resolution. It is the plural form: it reads "happen to have ...".
func (b *RepeatedlyToConditionBuilder[B, T]) HappenToHave(
	resolutions ...screenplay.Resolution[T],
) screenplay.Performable {
	return b.repeats(sentence.NewCondition(sentence.ShouldHave, true, resolutions))
}

// DoNotHappenTo tests that the answer satisfies none of the resolutions. It is the plural form: it reads "do not happen to ...".
func (b *RepeatedlyToConditionBuilder[B, T]) DoNotHappenTo(
	resolutions ...screenplay.Resolution[T],
) screenplay.Performable {
	return b.repeats(sentence.NewCondition(sentence.ShouldNot, true, resolutions))
}

// DoNotHappenToBe tests that the answer satisfies none of the resolutions. It is the plural form: it reads "do not happen to be ...".
func (b *RepeatedlyToConditionBuilder[B, T]) DoNotHappenToBe(
	resolutions ...screenplay.Resolution[T],
) screenplay.Performable {
	return b.repeats(sentence.NewCondition(sentence.ShouldNotBe, true, resolutions))
}

// DoNotHappenToHave tests that the answer satisfies none of the resolutions. It is the plural form: it reads "do not happen to have ...".
func (b *RepeatedlyToConditionBuilder[B, T]) DoNotHappenToHave(
	resolutions ...screenplay.Resolution[T],
) screenplay.Performable {
	return b.repeats(sentence.NewCondition(sentence.ShouldNotHave, true, resolutions))
}

// HappenToBeAnyOf tests that the answer satisfies at least one resolution. It is the plural form: it reads "happen to be any of: ...".
func (b *RepeatedlyToConditionBuilder[B, T]) HappenToBeAnyOf(
	resolutions ...screenplay.Resolution[T],
) screenplay.Performable {
	return b.repeats(sentence.NewCondition(sentence.ShouldAnyOf, true, resolutions))
}

// repeats creates the action repeating the body as long as the condition on
// the clause says so: while it holds, or until it holds.
func (b *RepeatedlyToConditionBuilder[B, T]) repeats(clause sentence.Clause[T]) *repeatedlyToAction[B] {
	term := sentence.NewTerm(b.question, clause)
	guard := condition.Term(term).IntroducedBy("while")
	if b.until {
		guard = guard.Negated().IntroducedBy("until")
	}

	return &repeatedlyToAction[B]{
		body:        b.body,
		atLeastOnce: b.atLeastOnce,
		condition:   &guard,
		mistake:     term.Mistake(),
	}
}

// RepeatedlyToCountBuilder holds the number of passes of a proceed.RepeatedlyTo
// action, and waits for Times.
type RepeatedlyToCountBuilder[B any] struct {
	body  B
	count int
}

// Times closes the number of passes. It reads "exactly 3 times".
func (b *RepeatedlyToCountBuilder[B]) Times() screenplay.Performable {
	return &repeatedlyToAction[B]{body: b.body, count: b.count}
}

// repeatedlyToAction is an action that performs its body again and again: a
// fixed number of times, or as long as its condition says so.
type repeatedlyToAction[B any] struct {
	body        B
	count       int
	atLeastOnce bool
	condition   *condition.Condition
	// mistake is why the condition can never be tested, when its sentence was
	// built wrongly.
	mistake error
}

// String describes the action: "proceed repeatedly to take a message while
// the queue length happens to be greater than 0".
func (a *repeatedlyToAction[B]) String() string {
	description := "proceed repeatedly to " + a.describeBody()
	if a.condition == nil {
		return fmt.Sprintf("%s exactly %s", description, timesOf(a.count))
	}

	if a.atLeastOnce {
		description += " at least once"
	}

	return fmt.Sprintf("%s %s %s", description, a.condition.Keyword, a.condition.Description)
}

// PerformAs performs the body pass after pass, and stops at the first pass the
// condition, or the count, does not allow. A failure of the body, or of the
// condition, ends the loop with an error that names the pass.
func (a *repeatedlyToAction[B]) PerformAs(theActor *screenplay.Actor) error {
	actionFor, mistake := a.actionsOfTheBody()
	if mistake == nil {
		mistake = a.mistakeOfTheClause()
	}

	if mistake != nil {
		return fmt.Errorf("%s failed to %s: %w", theActor.Name(), a, mistake)
	}

	for pass := 1; ; pass++ {
		goesOn, err := a.goesOnWith(theActor, pass)
		if err == nil && goesOn {
			err = theActor.AttemptsTo(actionFor(pass - 1))
		}

		if err != nil {
			return fmt.Errorf("%s failed to %s on the %s pass: %w", theActor.Name(), a, ordinalOf(pass), err)
		}

		if !goesOn {
			return nil
		}
	}
}

// goesOnWith reports whether the loop makes the given pass, counted from 1.
func (a *repeatedlyToAction[B]) goesOnWith(theActor *screenplay.Actor, pass int) (bool, error) {
	if a.condition == nil {
		return pass <= a.count, nil
	}

	if a.atLeastOnce && pass == 1 {
		return true, nil
	}

	return a.condition.HoldsFor(theActor)
}

// mistakeOfTheClause returns why the loop can never be performed, when its
// condition or its count was built wrongly.
func (a *repeatedlyToAction[B]) mistakeOfTheClause() error {
	if a.condition != nil {
		return a.mistake
	}

	if a.count < 0 {
		return fmt.Errorf("%w: the number of times %d is negative", screenplay.ErrInvalidResolution, a.count)
	}

	return nil
}

// actionsOfTheBody returns the function that gives the action of each pass,
// from its index: the body itself when it is an action, or the action the body
// returns when it is a function of the index.
func (a *repeatedlyToAction[B]) actionsOfTheBody() (func(index int) screenplay.Performable, error) {
	switch body := any(a.body).(type) {
	case screenplay.Performable:
		return func(int) screenplay.Performable { return body }, nil
	case func(index int) screenplay.Performable:
		if body != nil {
			return body, nil
		}
	}

	// A function that returns a concrete action, such as
	// func(int) *UploadThePartAction, does not match the type switch.
	function := reflect.ValueOf(a.body)
	if !isAFunctionOfTheIndex(function) {
		return nil, fmt.Errorf("%w: %T is neither an action nor a function of the index",
			screenplay.ErrInvalidResolution, a.body)
	}

	return func(index int) screenplay.Performable {
		argument := reflect.ValueOf(index).Convert(function.Type().In(0))
		action, _ := function.Call([]reflect.Value{argument})[0].Interface().(screenplay.Performable)

		return action
	}, nil
}

// isAFunctionOfTheIndex reports whether a value is a function that takes an
// integer and returns an action.
func isAFunctionOfTheIndex(function reflect.Value) bool {
	if function.Kind() != reflect.Func || function.IsNil() {
		return false
	}

	signature := function.Type()

	return signature.NumIn() == 1 && signature.In(0).Kind() == reflect.Int &&
		signature.NumOut() == 1 && signature.Out(0).Implements(reflect.TypeFor[screenplay.Performable]())
}

// describeBody describes the body: the description of the action, or the name
// of the function of the index written as words.
func (a *repeatedlyToAction[B]) describeBody() string {
	if action, isAnAction := any(a.body).(screenplay.Performable); isAnAction {
		return action.String()
	}

	function := reflect.ValueOf(a.body)
	if isAFunctionOfTheIndex(function) {
		if words := wordsOf(runtime.FuncForPC(function.Pointer()).Name()); words != "" {
			return words
		}
	}

	return "perform an action"
}

// wordsOf writes the name of a function as words: "pkg.UploadThePart" becomes
// "upload the part", and "pkg.SendHTTPRequest" "send HTTP request". It returns
// "" for a function literal, which has no name.
func wordsOf(qualifiedName string) string {
	name := strings.TrimSuffix(strings.ReplaceAll(qualifiedName, "[...]", ""), "-fm")
	name = name[strings.LastIndex(name, ".")+1:]

	if name == "" || unicode.IsDigit([]rune(name)[0]) || isAFunctionLiteral(name) {
		return ""
	}

	runes := []rune(name)
	words := []string{}
	start := 0

	for i := 1; i < len(runes); i++ {
		startsAWord := unicode.IsLower(runes[i-1]) || (i+1 < len(runes) && unicode.IsLower(runes[i+1]))
		if unicode.IsUpper(runes[i]) && startsAWord {
			words = append(words, string(runes[start:i]))
			start = i
		}
	}

	words = append(words, string(runes[start:]))

	for i, word := range words {
		// An acronym keeps its capitals: "HTTP" stays "HTTP".
		if strings.ToUpper(word) != word || len(word) == 1 {
			words[i] = strings.ToLower(word)
		}
	}

	return strings.Join(words, " ")
}

// isAFunctionLiteral reports whether the last part of the name of a function
// is the one Go gives a function literal: "func1", "func2", ...
func isAFunctionLiteral(name string) bool {
	digits, isNumbered := strings.CutPrefix(name, "func")

	return isNumbered && digits != "" && strings.Trim(digits, "0123456789") == ""
}

// timesOf counts the passes: "1 time", "3 times".
func timesOf(count int) string {
	if count == 1 {
		return "1 time"
	}

	return fmt.Sprintf("%d times", count)
}

// ordinalOf writes the rank of a pass: "1st", "2nd", "3rd", "4th", "11th", "21st".
func ordinalOf(pass int) string {
	rank := strconv.Itoa(pass)

	switch {
	case strings.HasSuffix(rank, "11"), strings.HasSuffix(rank, "12"), strings.HasSuffix(rank, "13"):
		return rank + "th"
	case strings.HasSuffix(rank, "1"):
		return rank + "st"
	case strings.HasSuffix(rank, "2"):
		return rank + "nd"
	case strings.HasSuffix(rank, "3"):
		return rank + "rd"
	default:
		return rank + "th"
	}
}

// repeatedlyToAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*repeatedlyToAction[screenplay.Performable])(nil)
