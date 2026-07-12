package action

import (
	"fmt"
	"strings"

	"github.com/grandper/go-screenplay/screenplay"
)

// ordinal renders 1 -> "first", 2 -> "second", ... for the boolean-branch description.
func ordinal(position int) string {
	names := []string{"first", "second", "third", "fourth", "fifth"}
	if position >= 1 && position <= len(names) {
		return names[position-1]
	}

	return fmt.Sprintf("%dth", position)
}

// joinChoices assembles "choose to A when x, to B when y, or to Z otherwise".
// Each clause is already in the infinitive form "to <action> when <condition>"; the
// fallback becomes the last alternative "to <action> otherwise", and the list is closed
// with "or" so it reads as picking exactly one branch (not a to-do list).
func joinChoices(clauses []string, fallback []screenplay.Performable) string {
	alternatives := make([]string, len(clauses), len(clauses)+1)
	copy(alternatives, clauses)
	if len(fallback) > 0 {
		alternatives = append(alternatives, "to "+performablesToString(fallback)+" otherwise")
	}

	if len(alternatives) <= 1 {
		return "choose " + strings.Join(alternatives, "")
	}

	last := len(alternatives) - 1

	return "choose " + strings.Join(alternatives[:last], ", ") + ", or " + alternatives[last]
}

// choiceBranch pairs a condition with the performables run when it is the first to hold.
// condition is the same type ConditionalAction uses, so When reuses its boolean logic and
// WhenThe reuses its question/resolution logic.
type choiceBranch struct {
	condition    condition
	performables []screenplay.Performable
}

// Choose starts a multi-way choice; each branch's condition is a boolean or a question/resolution.
func Choose() *ChoiceAction {
	return &ChoiceAction{
		branches: []choiceBranch{},
		pending:  []screenplay.Performable{},
		fallback: []screenplay.Performable{},
	}
}

// ChoiceAction performs the performables of the first branch whose condition holds.
type ChoiceAction struct {
	branches []choiceBranch
	pending  []screenplay.Performable
	fallback []screenplay.Performable
}

// To names the performables of the next branch; close it with When or WhenThe.
func (a *ChoiceAction) To(performables ...screenplay.Performable) *ChoiceAction {
	a.pending = performables

	return a
}

// When closes the pending branch with a boolean condition.
func (a *ChoiceAction) When(cond bool) *ChoiceAction {
	a.branches = append(a.branches, choiceBranch{
		condition: condition{
			describe: fmt.Sprintf("the %s condition holds", ordinal(len(a.branches)+1)),
			holds:    func(*screenplay.Actor) (bool, error) { return cond, nil },
		},
		performables: a.pending,
	})
	a.pending = []screenplay.Performable{}

	return a
}

// WhenThe closes the pending branch with a question/resolution condition.
func (a *ChoiceAction) WhenThe(
	question screenplay.Question,
	resolution screenplay.Resolution,
) *ChoiceAction {
	a.branches = append(a.branches, choiceBranch{
		condition:    questionCondition(question, resolution),
		performables: a.pending,
	})
	a.pending = []screenplay.Performable{}

	return a
}

// Otherwise closes the pending To branch as the fallback, run when no branch holds.
func (a *ChoiceAction) Otherwise() *ChoiceAction {
	a.fallback = a.pending
	a.pending = []screenplay.Performable{}

	return a
}

// Default closes the pending To branch as the fallback, run when no branch holds.
func (a *ChoiceAction) Default() *ChoiceAction {
	return a.Otherwise()
}

// String describes the action.
func (a *ChoiceAction) String() string {
	clauses := make([]string, 0, len(a.branches))
	for _, b := range a.branches {
		clauses = append(clauses, fmt.Sprintf(
			"to %s when %s", performablesToString(b.performables), b.condition.describe))
	}

	return joinChoices(clauses, a.fallback)
}

// PerformAs performs the task or the action as the provided actor.
func (a *ChoiceAction) PerformAs(theActor *screenplay.Actor) error {
	for _, b := range a.branches {
		holds, err := b.condition.holds(theActor)
		if err != nil {
			return err
		}

		if holds {
			return theActor.AttemptsTo(b.performables...)
		}
	}

	return theActor.AttemptsTo(a.fallback...)
}

// ChoiceAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*ChoiceAction)(nil)

// resolutionBranch pairs a resolution with the performables run when the answer first matches it.
type resolutionBranch struct {
	resolution   screenplay.Resolution
	performables []screenplay.Performable
}

// ChooseBasedOnThe starts a multi-way choice on the answer of a question; branches are resolutions.
func ChooseBasedOnThe(question screenplay.Question) *ChoiceOnAction {
	return &ChoiceOnAction{
		question: question,
		branches: []resolutionBranch{},
		pending:  []screenplay.Performable{},
		fallback: []screenplay.Performable{},
	}
}

// ChoiceOnAction performs the performables of the first branch whose resolution matches the answer.
type ChoiceOnAction struct {
	question screenplay.Question
	branches []resolutionBranch
	pending  []screenplay.Performable
	fallback []screenplay.Performable
}

// To names the performables of the next branch; close it with When.
func (a *ChoiceOnAction) To(performables ...screenplay.Performable) *ChoiceOnAction {
	a.pending = performables

	return a
}

// When closes the pending branch with the resolution the single answer must match.
func (a *ChoiceOnAction) When(resolution screenplay.Resolution) *ChoiceOnAction {
	a.branches = append(a.branches, resolutionBranch{resolution: resolution, performables: a.pending})
	a.pending = []screenplay.Performable{}

	return a
}

// Otherwise closes the pending To branch as the fallback, run when no branch matches.
func (a *ChoiceOnAction) Otherwise() *ChoiceOnAction {
	a.fallback = a.pending
	a.pending = []screenplay.Performable{}

	return a
}

// Default closes the pending To branch as the fallback, run when no branch matches.
func (a *ChoiceOnAction) Default() *ChoiceOnAction {
	return a.Otherwise()
}

// String describes the action.
func (a *ChoiceOnAction) String() string {
	clauses := make([]string, 0, len(a.branches))
	for i, b := range a.branches {
		subject := "it"
		if i == 0 {
			subject = "the " + a.question.String()
		}
		clauses = append(clauses, fmt.Sprintf(
			"to %s when %s is %s", performablesToString(b.performables), subject, b.resolution))
	}

	return joinChoices(clauses, a.fallback)
}

// PerformAs performs the task or the action as the provided actor.
func (a *ChoiceOnAction) PerformAs(theActor *screenplay.Actor) error {
	value, err := a.question.AnsweredBy(theActor)
	if err != nil {
		return err
	}

	for _, b := range a.branches {
		matched, matchErr := b.resolution.Resolve()(value)
		if matchErr != nil {
			return fmt.Errorf(
				"an error occurred when %s evaluated a branch on %s: %w",
				theActor.Name(), a.question, matchErr)
		}

		if matched {
			return theActor.AttemptsTo(b.performables...)
		}
	}

	return theActor.AttemptsTo(a.fallback...)
}

// ChoiceOnAction implements the screenplay.Performable interface.
var _ screenplay.Performable = (*ChoiceOnAction)(nil)
