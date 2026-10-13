package action

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/grandper/go-screenplay/v2/internal/reflected"
	"github.com/grandper/go-screenplay/v2/screenplay"
)

// PerformableFn represents is function that makes an actor performs an action or a task.
type PerformableFn func(theActor *screenplay.Actor) error

// FromFunc creates a new performable from a function.
func FromFunc(description string, fn PerformableFn) screenplay.Performable {
	return &funcPerformable{
		performAS:   fn,
		description: description,
	}
}

type funcPerformable struct {
	performAS   PerformableFn
	description string
}

// String describes the action.
func (p *funcPerformable) String() string {
	return p.description
}

// PerformAs performs the task or the action as the provided actor.
func (p *funcPerformable) PerformAs(actor *screenplay.Actor) error {
	if p.performAS == nil {
		return errors.New("the function given to FromFunc is nil")
	}

	return p.performAS(actor)
}

// funcPerformable implements the screenplay.Performable interface.
var _ screenplay.Performable = (*funcPerformable)(nil)

// Sentinel errors returned by FromFuncAndQuestions.
var (
	ErrTaskFuncNotFunction    = errors.New("taskFunc must be a function")
	ErrNoQuestions            = errors.New("at least one question must be provided")
	ErrFirstParamNotActor     = errors.New("first parameter must be *Actor")
	ErrQuestionCountMismatch  = errors.New("number of questions does not match function parameters")
	ErrTaskFuncMustReturnOne  = errors.New("taskFunc must return one value (Performable)")
	ErrTaskFuncMustReturnTask = errors.New("taskFunc must return a Performable")
	ErrAnswerNotAssignable    = errors.New("cannot assign or convert answer to expected type")
	ErrNotAQuestion           = errors.New("every question must be a screenplay.Question")
)

// FromFuncAndQuestions creates a performable from a description, a function and
// questions. When performed, the answers to the questions are asked to the actor,
// converted to the parameters of the function, and the returned task is performed.
//
// The questions may have different answer types, so they are received as plain
// values and checked when the action is performed: a value that is not a
// screenplay.Question fails with ErrNotAQuestion.
func FromFuncAndQuestions(
	description string,
	taskFunc any,
	questions ...any,
) screenplay.Performable {
	return &funcWithQuestionsPerformable{
		description: description,
		taskFunc:    taskFunc,
		questions:   questions,
	}
}

type funcWithQuestionsPerformable struct {
	description string
	taskFunc    any
	questions   []any
}

// String describes the action.
func (p *funcWithQuestionsPerformable) String() string {
	return p.description
}

// PerformAs performs the task or the action as the provided actor.
func (p *funcWithQuestionsPerformable) PerformAs(actor *screenplay.Actor) error { //nolint:gocognit
	funcValue := reflect.ValueOf(p.taskFunc)
	funcType := funcValue.Type()
	if funcType.Kind() != reflect.Func {
		return ErrTaskFuncNotFunction
	}
	if len(p.questions) == 0 {
		return ErrNoQuestions
	}
	if funcType.NumIn() < 1 || funcType.In(0) != reflect.TypeOf(&screenplay.Actor{}) {
		return ErrFirstParamNotActor
	}
	if funcType.NumIn()-1 != len(p.questions) {
		return ErrQuestionCountMismatch
	}
	if funcType.NumOut() != 1 {
		return ErrTaskFuncMustReturnOne
	}
	taskType := reflect.TypeFor[screenplay.Performable]()
	if !funcType.Out(0).Implements(taskType) {
		return ErrTaskFuncMustReturnTask
	}
	args := []reflect.Value{reflect.ValueOf(actor)}
	for i, value := range p.questions {
		q, isAQuestion := reflected.QuestionOf(value, funcType.In(0))
		if !isAQuestion {
			return fmt.Errorf("%w: got a %T at position %d", ErrNotAQuestion, value, i)
		}
		answer, err := q.AnsweredBy(actor)
		if err != nil {
			return err
		}
		expectedType := funcType.In(i + 1)
		answerValue := reflect.ValueOf(answer)
		if !answerValue.IsValid() {
			return ErrAnswerNotAssignable
		}
		if !answerValue.Type().AssignableTo(expectedType) {
			if answerValue.Type().ConvertibleTo(expectedType) {
				answerValue = answerValue.Convert(expectedType)
			} else {
				return ErrAnswerNotAssignable
			}
		}
		args = append(args, answerValue)
	}
	results := funcValue.Call(args)
	task, ok := results[0].Interface().(screenplay.Performable)
	if !ok {
		return ErrTaskFuncMustReturnTask
	}
	return task.PerformAs(actor)
}

// funcWithQuestionsPerformable implements the screenplay.Performable interface.
var _ screenplay.Performable = (*funcWithQuestionsPerformable)(nil)
